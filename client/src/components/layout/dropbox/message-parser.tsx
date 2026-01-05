import { Fragment } from 'react';

type NodeProcessor = (node: Node, index: number) => React.ReactNode[];

export default function parseMessage(message: string) {
    const parser = new DOMParser();
    const doc = parser.parseFromString(message, 'text/html');

    const processTextNode: NodeProcessor = (node, index) => {
        const lines = node.textContent?.split('\n') || [];
        return lines.map((line, i) => (
            <Fragment key={`text-${index}-${i}`}>
                {line}
                {i < lines.length - 1 && <br />}
            </Fragment>
        ));
    };

    const processSpanElement: NodeProcessor = (node, index) => {
        const element = node as Element;
        const lines = element.innerHTML.split('\n') || [];
        return lines.map((line, i) => (
            <Fragment key={`span-${index}-${i}`}>
                <span className={element.className}>{line}</span>
                {i < lines.length - 1 && <br />}
            </Fragment>
        ));
    };

    const processNode = (node: Node, index: number): React.ReactNode[] => {
        if (node.nodeType === Node.TEXT_NODE) {
            return processTextNode(node, index);
        }

        if (node.nodeType === Node.ELEMENT_NODE && (node as Element).tagName.toLowerCase() === 'span') {
            return processSpanElement(node, index);
        }

        return [];
    };

    const elements = Array.from(doc.body.childNodes).map(processNode);
    return elements;
}
