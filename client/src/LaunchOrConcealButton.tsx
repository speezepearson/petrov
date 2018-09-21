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
                    {props.impactTime ? <Timer currentTime={props.currentTime} zeroTime={props.impactTime} showHours={false} /> : ''}
                </div>
                <div style={{alignSelf: 'flex-end'}}>
                    (click to hide)
                </div>
            </div>
        </div>
    );
}
