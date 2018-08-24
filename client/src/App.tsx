import jQuery from 'jquery';
import React from 'react';

function updateAppForever(app: App, playerName: string): () => void {
    let stopped = false;

    function loop() {
        if (stopped) return;
        jQuery.get(
            `/${playerName}`,
            function(data) {
                setTimeout(loop, 1000);
                app.setState(data);
            }
        );
    }

    loop();

    return () => {stopped = true;};
}

type AppProps = {
    playerName: string;
}
type AppState = {
    Phase: string;
    TimeRemainingNs: number;
    AlarmTimesNsRemaining: number[];
    KilledBy: string;
    TimeToMyImpactNs?: number;
}

export class App extends React.Component<AppProps, AppState> {
    private stopUpdating?: () => void;
    componentWillMount() {
        this.stopUpdating = updateAppForever(this, this.props.playerName);
    }
    componentWillUnmount() {
        if (this.stopUpdating) {
            this.stopUpdating();
        }
    }
    render() {
        return <div>&#9762; Goodbye world! &#9762;</div>;
    }
}
