import jQuery from 'jquery';
import React from 'react';

type LaunchButtonProps = {
    playerName: string;
    alreadyLaunched: boolean;
}
type LaunchButtonState = {
    alreadyLaunched: boolean;
}

export class LaunchButton extends React.Component<LaunchButtonProps, LaunchButtonState> {
    constructor(props: LaunchButtonProps) {
        super(props);
        this.state = {
            alreadyLaunched: props.alreadyLaunched,
        };
    }
    render() {
        return <button disabled={this.state.alreadyLaunched} onClick={() => this.launch()}>
            Launch
        </button>;
    }

    launch() {
        this.setState({alreadyLaunched: true});
        jQuery.post({
            url: `/${this.props.playerName}/launch`,
        });
    }
}