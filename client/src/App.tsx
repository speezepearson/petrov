import jQuery from 'jquery';
import React from 'react';

import { Timer } from './Timer'
import {LaunchOrConcealButton} from "./LaunchOrConcealButton";
import './App.css';

enum Phase {
    ENDED = "Ended",
    PRESTART = "PreStart",
    RUNNING = "Running",
    OVERTIME = "Overtime",
}

function nowPlus(seconds: number): Date {
    const result = new Date();
    result.setSeconds(result.getSeconds() + seconds);
    return result;
}

type AppProps = {
    playerName: string;
}
type AppState = {
    phase?: string;
    timeRemaining: number;
    alarmTimesRemaining: number[];
    killedBy: string;
    timeToMyImpact?: number;
}

export class App extends React.Component<AppProps, AppState> {
    private updaterId?: number;
    constructor(props: AppProps) {
      super(props);
      this.state = {
        phase: undefined,
        timeRemaining: 0,
        alarmTimesRemaining: [],
        killedBy: '',
        timeToMyImpact: undefined,
      };
    }
    componentWillMount() {
        this.updaterId = setInterval(this.fetchData.bind(this), 1000);
    }
    componentWillUnmount() {
        if (this.updaterId) {
            clearInterval(this.updaterId);
        }
    }
    render() {
        console.log("Rendering:", this.state);
        switch (this.state.phase) {

            case undefined:
                return 'Loading...';

            case Phase.PRESTART:
                return 'Game not yet started.';

            case Phase.RUNNING:
            case Phase.OVERTIME:
                let incoming: boolean = (this.state.alarmTimesRemaining.length > 0);
                return <div>
                    <div id="time-remaining">
                        <Timer zeroTime={nowPlus(this.state.timeRemaining)} /> remaining
                    </div>

                    <div id="top-stuff">
                        {
                            incoming
                            ? [
                                <div style={{flexGrow: 1}} key="INCOMING">INCOMING</div>,
                                ...this.state.alarmTimesRemaining.map((d, i) => (
                                        <div style={{flexGrow: 1}} key={i}>
                                            <Timer zeroTime={nowPlus(d)} />
                                        </div>
                                    )),
                                <div style={{flexGrow: 1}} key="LAUNCH NOW">LAUNCH NOW</div>
                              ]
                            : ''
                        }
                    </div>

                    <div id="bottom-stuff">
                        <LaunchOrConcealButton
                            playerName={this.props.playerName}
                            impactTime={this.state.timeToMyImpact ? nowPlus(this.state.timeToMyImpact) : null}
                        />
                    </div>
                </div>;

            case Phase.ENDED:
                return `Game over! You're alive! Everyone else is ${this.state.timeToMyImpact ? "dead. Remember? You killed them." : "alive too!"}`;

            default:
                return `Unknown phase: ${this.state.phase}`;
        }
    }

    fetchData() {
        jQuery.get({
            url: `/${this.props.playerName}/status`,
            success: dataText => {
                const data = JSON.parse(dataText);
                console.log(data);
                this.setState({
                    phase: data.Phase,
                    timeRemaining: data.TimeRemaining / 1e9,
                    alarmTimesRemaining: (data.AlarmTimesRemaining || []).map((x:number) => x/1e9),
                    killedBy: data.KilledBy,
                    timeToMyImpact: data.TimeToMyImpact || null,
                });
            },
        });
    }
}
