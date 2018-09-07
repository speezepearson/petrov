import jQuery from 'jquery';
import React from 'react';
import { Timer } from './Timer';

import './LaunchOrConcealButton.css';

type LaunchOrConcealButtonProps = {
    playerName: string;
    impactTime: Date | null;
    currentTime: Date;
}

export function LaunchOrConcealButton(props: LaunchOrConcealButtonProps) {
    return (
        <button className={`launch-button launch-button--${props.impactTime ? 'ticking' : 'ready'}`}
                onClick={() => jQuery.post({
                    url: `/${props.playerName}/${(!!props.impactTime) ? 'conceal' : 'launch'}`,
                })}>
            {props.impactTime ? <Timer currentTime={props.currentTime} zeroTime={props.impactTime}/> : ''}
            {props.impactTime ? <br /> : ''}
            {props.impactTime ? "Feign innocence" : "Launch"}
        </button>
    );
}
