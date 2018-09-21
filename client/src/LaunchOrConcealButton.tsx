import jQuery from 'jquery';
import React from 'react';
import { Timer } from './Timer';

import './LaunchOrConcealButton.css';

type LaunchOrConcealButtonProps = {
    impactTime: Date | null;
    currentTime: Date;
    onLaunch: () => void;
    onConceal: () => void;
}

function formatTimeToImpact(impactTime: Date | null, currentTime: Date) {
    if (!impactTime) return '(unlaunched)';
    if (impactTime < currentTime) return '(landed)';
    return <Timer currentTime={currentTime} zeroTime={impactTime} showHours={false} />
}

export function LaunchOrConcealButton(props: LaunchOrConcealButtonProps) {
    return (
        <div className="launch-button-container">
            <button className={`launch-button launch-button--${props.impactTime ? 'ticking' : 'ready'}`}
                    onClick={props.onLaunch}>
                Launch
            </button>
            <div className="time-to-impact" style={{display: 'flex', flexDirection: 'row'}}>
                <div style={{alignSelf: 'flex-start'}}>
                    Time to impact:
                </div>
                <div onClick={props.onConceal} style={{flexGrow: 1, cursor: 'pointer'}}>
                    {formatTimeToImpact(props.impactTime, props.currentTime)}
                </div>
                <div style={{alignSelf: 'flex-end'}}>
                    (click to hide)
                </div>
            </div>
        </div>
    );
}
