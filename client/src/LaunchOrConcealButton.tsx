import jQuery from 'jquery';
import React from 'react';

import './LaunchOrConcealButton.css';

type LaunchOrConcealButtonProps = {
    playerName: string;
    launched: boolean;
}
type LaunchOrConcealButtonState = {
    launched: boolean;
}

export class LaunchOrConcealButton extends React.Component<LaunchOrConcealButtonProps, LaunchOrConcealButtonState> {
    constructor(props: LaunchOrConcealButtonProps) {
        super(props);
        this.state = {
            launched: props.launched,
        };
    }
    render() {
        return (
            <button className={`launch-button launch-button--${this.state.launched ? 'ticking' : 'ready'}`}
                    onClick={() => this.launchOrConceal()}>
                {this.state.launched ? "Feign innocence" : "Launch"}
            </button>
        );
    }

    launchOrConceal() {
        const hadLaunched: boolean = this.state.launched;
        this.setState({launched: !hadLaunched});
        jQuery.post({
            url: `/${this.props.playerName}/${hadLaunched ? "conceal" : "launch"}`,
        });
    }
}
