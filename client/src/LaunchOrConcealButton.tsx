import jQuery from 'jquery';
import React from 'react';
import { Timer } from './Timer';

import './LaunchOrConcealButton.css';

type LaunchOrConcealButtonProps = {
    playerName: string;
    impactTime: Date | null;
    currentTime: Date;
    onClick: () => void;
}

export function LaunchOrConcealButton(props: LaunchOrConcealButtonProps) {
    return (
        <button className={`launch-button launch-button--${props.impactTime ? 'ticking' : 'ready'}`}
                onClick={props.onClick}>
            {props.impactTime ? <Timer currentTime={props.currentTime} zeroTime={props.impactTime}/> : ''}
            {props.impactTime ? <br /> : ''}
            {props.impactTime ? "Feign innocence" : "Launch"}
        </button>
    );
}
