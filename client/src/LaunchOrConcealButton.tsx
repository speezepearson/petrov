import jQuery from 'jquery';
import React from 'react';
import { Timer } from './Timer';

import './LaunchOrConcealButton.css';

const EXPECTED_MISSILE_FLIGHT_TIME = 10;

type LaunchOrConcealButtonProps = {
    playerName: string;
    impactTime: Date | null;
}
type LaunchOrConcealButtonState = {
    impactTime: Date | null;
}

export class LaunchOrConcealButton extends React.Component<LaunchOrConcealButtonProps, LaunchOrConcealButtonState> {
    constructor(props: LaunchOrConcealButtonProps) {
        super(props);
        this.state = {
            impactTime: props.impactTime,
        };
    }
    render() {
        return (
            <button className={`launch-button launch-button--${this.state.impactTime ? 'ticking' : 'ready'}`}
                    onClick={() => this.launchOrConceal()}>
                {
                    this.state.impactTime
                    ? [<Timer zeroTime={this.state.impactTime}/>, <br />]
                    : ''
                }
                {this.state.impactTime ? "Feign innocence" : "Launch"}
            </button>
        );
    }

    launchOrConceal() {
        const hadLaunched: boolean = !!this.state.impactTime;
        if (hadLaunched) {
            this.setState({impactTime: null});
            jQuery.post({
                url: `/${this.props.playerName}/conceal`,
            });
        } else {
            const impactTime = new Date();
            impactTime.setSeconds(impactTime.getSeconds() + EXPECTED_MISSILE_FLIGHT_TIME);
            this.setState({impactTime: impactTime});
            jQuery.post({
                url: `/${this.props.playerName}/launch`,
            });
        }
    }
}
