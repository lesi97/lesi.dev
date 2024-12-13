'use client';

import React, { useEffect, useState } from 'react';
import Snowfall from 'react-snowfall';
import { StaticImageData } from 'next/image';

import heartImg1 from '../../_assets/_images/heart-1.webp';
import heartImg2 from '../../_assets/_images/heart-2.webp';
import heartImg3 from '../../_assets/_images/heart-3.webp';

import pumpkinImg1 from '../../_assets/_images/pumpkin-1.webp';
import spiderWebImg1 from '../../_assets/_images/spider-web-1.webp';

import snowflakeImg1 from '../../_assets/_images/snowflake-1.webp';
import snowflakeImg2 from '../../_assets/_images/snowflake-2.webp';
import snowflakeImg3 from '../../_assets/_images/snowflake-3.webp';

import fireworkImg1 from '../../_assets/_images/firework-1.webp';

const SnowfallWrapper: React.FC<{ season: string | null }> = ({ season }) => {
    const [snowfallData, setSnowfallData] = useState({
        snowflakeCount: 150,
        radius: [0.25, 15],
        rotationSpeed: [-1, 1],
        opacity: [1, 1],
        speed: [1, 3],
        images: [],
    });

    useEffect(() => {
        const valentines = [heartImg1, heartImg2, heartImg3];
        const halloween = [pumpkinImg1, spiderWebImg1];
        const christmas = [snowflakeImg1, snowflakeImg2, snowflakeImg3];
        const newYears = [fireworkImg1];

        // const testSeason = 'Halloween';

        // https://github.com/cahilfoley/react-snowfall
        switch (season) {
            case 'Valentines':
                setSnowfallData({
                    snowflakeCount: 150,
                    radius: [0.25, 15],
                    rotationSpeed: [-0.5, 0.5],
                    opacity: [0.2, 0.6],
                    speed: [0.1, 1],
                    images: loadImages(valentines),
                });
                break;
            case 'Halloween':
                setSnowfallData({
                    snowflakeCount: 80,
                    radius: [5, 35],
                    rotationSpeed: [-1, 1],
                    opacity: [1, 1],
                    speed: [1, 3],
                    images: loadImages(halloween),
                });
                break;
            case 'Lesi-Birthday':
                setSnowfallData({
                    snowflakeCount: 150,
                    radius: [0.25, 15],
                    rotationSpeed: [-1, 1],
                    opacity: [1, 1],
                    speed: [1, 3],
                    images: loadImages(halloween),
                });
                break;
            case 'Christmas':
                setSnowfallData({
                    snowflakeCount: 150,
                    radius: [0.25, 15],
                    rotationSpeed: [-1, 1],
                    opacity: [1, 1],
                    speed: [1, 3],
                    images: loadImages(christmas),
                });
                break;
            case 'New-Years':
                setSnowfallData({
                    snowflakeCount: 100,
                    radius: [10, 20],
                    rotationSpeed: [-2, 2],
                    opacity: [1, 1],
                    speed: [0.5, 2],
                    images: loadImages(newYears),
                });
                break;
            default:
                setSnowfallData({
                    snowflakeCount: null,
                    radius: [null, null],
                    rotationSpeed: [null, null],
                    opacity: [null, null],
                    speed: [null, null],
                    images: [null],
                });
                break;
        }

        function loadImages(imageType: StaticImageData[]) {
            const staticImages = imageType;
            const loadedImages = staticImages.map((imgData) => {
                const img = new Image();
                img.src = imgData.src;
                return img;
            });
            return loadedImages;
        }
    }, []);

    if (snowfallData.images.length === 0) return null;

    return (
        <div
            style={{
                position: 'fixed',
                top: 0,
                left: 0,
                width: '100vw',
                height: '100vh',
                pointerEvents: 'none',
                zIndex: 0,
            }}
        >
            <Snowfall
                snowflakeCount={snowfallData.snowflakeCount}
                images={snowfallData.images}
                radius={snowfallData.radius}
                rotationSpeed={snowfallData.rotationSpeed}
                opacity={snowfallData.opacity}
                speed={snowfallData.speed}
                style={{
                    top: '47px',
                    height: 'calc(100vh - 47px)',
                    position: 'fixed',
                    pointerEvents: 'none',
                }}
            />
        </div>
    );
};

export default SnowfallWrapper;
