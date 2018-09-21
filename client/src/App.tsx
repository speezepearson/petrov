import jQuery from 'jquery';
import React from 'react';

import { Timer } from './Timer'
import { LaunchOrConcealButton } from "./LaunchOrConcealButton";
import { KLAXON, FWOOSH, KABOOM } from './sounds';
import './App.css';

enum Phase {
    ENDED = "Ended",
    PRESTART = "PreStart",
    RUNNING = "Running",
}

function nowPlus(seconds: number): Date {
    const result = new Date();
    result.setSeconds(result.getSeconds() + seconds);
    return result;
}

type AppProps = {
    password: string;
}
type AppState = {
    interactedWith?: boolean;
    phase?: string;
    gameEndTime?: Date;
    alarmImpactTimes: Date[];
    killedBy: string;
    myImpactTime?: Date;
    currentTime: Date;
    lastSynced?: Date;
    playerName?: string;
    missileFlightTime?: number;
}

export class App extends React.Component<AppProps, AppState> {
    private fetcherId?: number;
    private tickerId?: number;
    constructor(props: AppProps) {
      super(props);
      this.state = {
        phase: undefined,
        gameEndTime: undefined,
        alarmImpactTimes: [],
        killedBy: '',
        myImpactTime: undefined,
        currentTime: new Date(),
        lastSynced: undefined,
        playerName: undefined,
      };
    }
    componentWillMount() {
        this.fetcherId = window.setInterval(() => this.fetchData(), 1000);
        this.tickerId = window.setInterval(() => this.setState({currentTime: new Date()}), 10);
    }
    componentWillUnmount() {
        if (this.fetcherId) {
            clearInterval(this.fetcherId);
        }
        if (this.tickerId) {
            clearInterval(this.tickerId);
        }
    }

    isSyncFailing(): boolean {
        return (!!this.state.lastSynced && (this.state.lastSynced < nowPlus(-5)))
    }

    render() {
        if (!this.state.interactedWith) {
            return <div id="modal">
                <div id="modal__content">
                    <button className="connect-button" onClick={() => this.setState({interactedWith: true})}>
                        Connect
                    </button>
                </div>
            </div>
        }
        if (this.isSyncFailing()) {
            KLAXON.play();
            return <div id="modal" style={{backgroundColor: 'green'}}>
                <div id="modal__content">
                    Lost connection to server! <br />
                    NOTIFY GAME CONTROL.
                </div>
            </div>;
        }
        return [this.renderOverlay(), this.renderGame()];
    }

    renderOverlay() {
        const timeRemaining = (
            (this.state.phase === Phase.RUNNING)
            ? <div key="time-remaining" id="time-remaining">
                {this.state.gameEndTime ? <div><Timer currentTime={this.state.currentTime} zeroTime={this.state.gameEndTime} showHours={true} /> remaining</div> : ''}
              </div>
            : ''
        );
        return [
            timeRemaining,
            <div key="player-name" id="player-name">
                {this.state.playerName}
            </div>,
        ];
    }

    renderGame() {

        (window as any)._sounds = {KLAXON, FWOOSH, KABOOM};

        if (this.state.killedBy.length > 0) {
            return <div id="modal">
                <div id="modal__content">
                    You were killed by {this.state.killedBy}.
                </div>
            </div>
        }

        if (this.state.missileFlightTime && (this.state.alarmImpactTimes.filter(x => (x > nowPlus(this.state.missileFlightTime!-5))).length > 0)) {
            KLAXON.play();
        }

        switch (this.state.phase) {

            case undefined:
                return <div id="modal">
                    <div id="modal__content">
                        Connecting...
                    </div>
                </div>;

            case Phase.PRESTART:
                return <div id="modal">
                    <div id="modal__content">
                        Waiting for early-warning system to come online...
                    </div>
                </div>;

            case Phase.RUNNING:
                const impactTimes: Date[] = this.state.alarmImpactTimes.filter(t => t > this.state.currentTime);
                const incoming: boolean = (impactTimes.length > 0);
                return [

                    <div key="top-stuff" id="top-stuff">
                        {
                            incoming
                            ? [
                                <div id="incoming-label" key="INCOMING">
                                    {impactTimes.length == 1
                                     ? "MISSILE INCOMING"
                                     : `${impactTimes.length} MISSILES INCOMING`}
                                </div>,
                                <div id="incoming-timers" key="incoming-timers">
                                    {impactTimes.sort().map((d, i) => (
                                        <div className="incoming-timers__timer-wrapper" key={i}>
                                            <div className={`incoming-timers__timer ${i===0 ? 'large' : ''}`}>
                                                <Timer currentTime={this.state.currentTime} zeroTime={d} />
                                            </div>
                                        </div>
                                    ))}
                                </div>,
                                <div id="launch-now-label" key="LAUNCH NOW">LAUNCH NOW</div>
                              ]
                            : ''
                        }
                    </div>,

                    <div key="bottom-stuff" id="bottom-stuff">
                        <LaunchOrConcealButton
                            impactTime={this.state.myImpactTime || null}
                            onConceal={() => {
                                if (!this.state.myImpactTime) return;
                                this.setState({
                                    myImpactTime: undefined,
                                });
                                jQuery.post({
                                    url: `/${this.props.password}/conceal`,
                                })
                            }}
                            onLaunch={() => {
                                if (!!this.state.myImpactTime) return;
                                this.setState({
                                    myImpactTime: nowPlus(this.state.missileFlightTime!),
                                });
                                jQuery.post({
                                    url: `/${this.props.password}/launch`,
                                })
                            }}
                            currentTime={this.state.currentTime}
                        />
                    </div>,
                ];

            case Phase.ENDED:
                return <div id="modal">
                    <div id="modal__content">
                        Game over! You're alive! Everyone else is {this.state.myImpactTime ? ["dead.", <br/>, "Remember? You killed them."] : "alive too!"}
                    </div>
                </div>;

            default:
                return <div id="modal">
                    <div id="modal__content">
                        Unknown phase: ${this.state.phase}
                    </div>
                </div>;
        }
    }

    fetchData() {
        jQuery.get({
            url: `/${this.props.password}/status`,
            success: dataText => {
                const data = JSON.parse(dataText);
                const now = new Date();
                console.log("received", data, "at", now);
                if (data.KilledBy && !this.state.killedBy) {
                    console.log('playing kaboom');
                    KLAXON.reset();
                    FWOOSH.reset();
                    KABOOM.play();
                }
                this.setState({
                    phase: data.Phase,
                    gameEndTime: nowPlus(data.TimeRemaining / 1e9),
                    alarmImpactTimes: (data.AlarmTimesRemaining || []).map((x:number) => nowPlus(x/1e9)),
                    killedBy: data.KilledBy,
                    myImpactTime: data.TimeToMyImpact ? nowPlus(data.TimeToMyImpact / 1e9) : undefined,
                    lastSynced: now,
                    playerName: data.Player,
                    missileFlightTime: data.MissileFlightTime / 1e9,
                });
            },
        });
    }
}
