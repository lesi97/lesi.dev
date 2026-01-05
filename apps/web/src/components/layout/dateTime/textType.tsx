import { ElementType, useEffect, useRef, useState, createElement, useMemo, useCallback } from 'react';
import { gsap } from 'gsap';
import { mergeClassNames } from '@/utils';
import { ScrambledText } from './scrambleText';

interface TextTypeProps {
    className?: string;
    showCursor?: boolean;
    hideCursorWhileTyping?: boolean;
    cursorCharacter?: string | React.ReactNode;
    cursorBlinkDuration?: number;
    cursorClassName?: string;
    text: string | string[];
    as?: ElementType;
    typingSpeed?: number;
    initialDelay?: number;
    pauseDuration?: number;
    deletingSpeed?: number;
    loop?: boolean;
    variableSpeed?: { min: number; max: number };
    onSentenceComplete?: (sentence: string, index: number) => void;
    startOnVisible?: boolean;
    reverseMode?: boolean;
    onTypingComplete?: () => void;
}

export function TextType({
    text,
    as: Component = 'div',
    typingSpeed = 50,
    initialDelay = 0,
    pauseDuration = 2000,
    deletingSpeed = 30,
    loop = false,
    className = '',
    showCursor = true,
    hideCursorWhileTyping = false,
    cursorCharacter = '|',
    cursorClassName = '',
    cursorBlinkDuration = 0.5,
    variableSpeed,
    onSentenceComplete,
    startOnVisible = false,
    reverseMode = false,
    onTypingComplete,
    ...props
}: TextTypeProps & React.HTMLAttributes<HTMLElement>) {
    const [displayedText, setDisplayedText] = useState('');
    const [currentCharIndex, setCurrentCharIndex] = useState(0);
    const [isDeleting, setIsDeleting] = useState(false);
    const [currentTextIndex, setCurrentTextIndex] = useState(0);
    const [isVisible, setIsVisible] = useState(!startOnVisible);
    const [hasCompletedOnce, setHasCompletedOnce] = useState(false);
    const cursorRef = useRef<HTMLSpanElement>(null);
    const containerRef = useRef<HTMLElement>(null);

    const textArray = useMemo(() => (Array.isArray(text) ? text : [text]), [text]);

    const getRandomSpeed = useCallback(() => {
        if (!variableSpeed) {
            return typingSpeed;
        }
        const { min, max } = variableSpeed;
        return Math.random() * (max - min) + min;
    }, [variableSpeed, typingSpeed]);

    useEffect(() => {
        if (!startOnVisible || !containerRef.current) {
            return;
        }
        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) {
                        setIsVisible(true);
                    }
                });
            },
            { threshold: 0.1 }
        );
        observer.observe(containerRef.current);
        return function cleanup() {
            observer.disconnect();
        };
    }, [startOnVisible]);

    useEffect(() => {
        if (showCursor && cursorRef.current) {
            gsap.set(cursorRef.current, { opacity: 1 });
            gsap.to(cursorRef.current, {
                opacity: 0,
                duration: cursorBlinkDuration,
                repeat: -1,
                yoyo: true,
                ease: 'power2.inOut',
            });
        }
    }, [showCursor, cursorBlinkDuration]);

    useEffect(() => {
        if (!hasCompletedOnce) {
            return;
        }
        const next = reverseMode ? textArray[0].split('').reverse().join('') : textArray[0];
        if (displayedText !== next) {
            setDisplayedText(next);
            setCurrentCharIndex(next.length);
        }
    }, [hasCompletedOnce, textArray, reverseMode, displayedText]);

    useEffect(() => {
        if (!isVisible) {
            return;
        }
        let timeout: number;
        const currentText = textArray[currentTextIndex];
        const processedText = reverseMode ? currentText.split('').reverse().join('') : currentText;

        const executeTypingAnimation = () => {
            if (isDeleting) {
                if (displayedText === '') {
                    setIsDeleting(false);
                    if (currentTextIndex === textArray.length - 1 && !loop) {
                        setHasCompletedOnce(true);
                        if (onSentenceComplete) {
                            onSentenceComplete(textArray[currentTextIndex], currentTextIndex);
                        }
                        if (onTypingComplete) {
                            onTypingComplete();
                        }
                        return;
                    }
                    if (onSentenceComplete) {
                        onSentenceComplete(textArray[currentTextIndex], currentTextIndex);
                    }
                    setCurrentTextIndex((prev) => (prev + 1) % textArray.length);
                    setCurrentCharIndex(0);
                    timeout = window.setTimeout(() => {
                        setDisplayedText('');
                    }, pauseDuration);
                } else {
                    timeout = window.setTimeout(() => {
                        setDisplayedText((prev) => prev.slice(0, -1));
                    }, deletingSpeed);
                }
            } else {
                if (currentCharIndex < processedText.length) {
                    timeout = window.setTimeout(
                        () => {
                            setDisplayedText((prev) => prev + processedText[currentCharIndex]);
                            setCurrentCharIndex((prev) => prev + 1);
                            if (currentCharIndex + 1 === processedText.length && (!loop || textArray.length === 1)) {
                                setHasCompletedOnce(true);
                                if (onTypingComplete) {
                                    onTypingComplete();
                                }
                            }
                        },
                        variableSpeed ? getRandomSpeed() : typingSpeed
                    );
                } else if (textArray.length > 1) {
                    timeout = window.setTimeout(() => {
                        setIsDeleting(true);
                    }, pauseDuration);
                } else if (!loop) {
                    setHasCompletedOnce(true);
                    if (onTypingComplete) {
                        onTypingComplete();
                    }
                }
            }
        };

        if (currentCharIndex === 0 && !isDeleting && displayedText === '') {
            timeout = window.setTimeout(executeTypingAnimation, initialDelay);
        } else {
            executeTypingAnimation();
        }

        return function cleanup() {
            window.clearTimeout(timeout);
        };
    }, [
        currentCharIndex,
        displayedText,
        isDeleting,
        typingSpeed,
        deletingSpeed,
        pauseDuration,
        textArray,
        currentTextIndex,
        loop,
        initialDelay,
        isVisible,
        reverseMode,
        variableSpeed,
        onSentenceComplete,
        getRandomSpeed,
        onTypingComplete,
    ]);

    const shouldHideCursor =
        hideCursorWhileTyping && (currentCharIndex < textArray[currentTextIndex].length || isDeleting);

    const useScramble = hasCompletedOnce;

    return createElement(
        Component,
        {
            ref: containerRef,
            className: mergeClassNames('tracking-tight inline-flex items-baseline', className),
            ...props,
        },
        useScramble ? (
            <ScrambledText active radius={100} duration={1.2} speed={0.5} className='text-primary'>
                {displayedText}
            </ScrambledText>
        ) : (
            <span className='text-primary'>{displayedText}</span>
        ),
        showCursor && (
            <span
                ref={cursorRef}
                className={mergeClassNames('ml-1 opacity-100', shouldHideCursor ? 'hidden' : '', cursorClassName)}>
                {cursorCharacter}
            </span>
        )
    );
}
