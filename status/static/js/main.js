const domContainer = document.querySelector('#app');

class GlobalStatus extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        let text = 'Todos Servicios Operativos.';
        let c = 'alert alert-success';

        for (let e in this.props.pages) {
            let service = this.props.pages[e];

            if (!service.up) {
                c = 'alert alert-danger';
                text = 'Servicios Temporalmente Cortados.';
                break;
            }
        }

        return (
            <div className={c} role="alert">
                {text}
            </div>
        );
    }
}

class IncidentPopup extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        if (!this.props.popup) return <span></span>;

        return (
            <div className="status-popup-container" style={{display: this.props.popup.show ? "block" : "none"}}>
                <div className="status-popup">

                </div>
            </div>
        );
    }
}

class PastIncidents extends React.Component {
    constructor(props) {
        super(props);
    }

    _buildIncidents() {
        let divs = [];

        const h = {};

        for (let url in this.props.history) {
            h[url] = this.props.history[url].slice().reverse();
        }

        let len = 30;

        let historyIndices = {};

        for (let i = 0; i < len; i++) {
            let day = new Date();
            day.setDate(day.getDate() - i);
            day.setHours(0, 0, 0, 0);

            // Subtract days
            let unix = day.getTime();
            let incidents = [];
            
            for (let url in h) {
                if (!(url in this.props.pages)) continue;
                let index = url in historyIndices ? historyIndices[url] : 0;
                let prevStatus = null;
                let prevTime = null;

                while (index < h[url].length &&
                    h[url][index].unix > unix)
                {
                    if (prevStatus == null) prevStatus = h[url][index].up;
                    if (prevTime == null) prevTime = h[url][index].timestamp;
                    if (h[url][index].up ^ prevStatus) {
                        let time = new Date(prevTime).toLocaleTimeString();
                        let text = prevStatus ? "Ha vuelto" : "Ha caído";
                        let name = this.props.pages[url].name;
                        incidents.unshift({
                            "time": prevTime,
                            "el": <li>{time + ': ' + text + ' ' + name}</li>
                        });
                    }
                    
                    prevStatus = h[url][index].up;
                    prevTime = h[url][index].timestamp;
                    index++;
                }

                historyIndices[url] = index;
            }

            let daytext = <p className="text-muted">No se reportaron incidencias.</p>;

            if (incidents.length) {
                incidents.sort((a, b) => new Date(a.time) - new Date(b.time));
                daytext = <ul>{incidents.map((v) => v.el)}</ul>;
            }

            divs.push(
                <div key={i}>
                    <h4 className="pt-3">{day.toLocaleDateString()}</h4>
                    <hr className="mt-0" />
                    {daytext}
                </div>
            );
        }

        return divs;
    }

    render() {
        if (!this.props.history || !this.props.pages) {
            return <span></span>
        }

        const incidents = this._buildIncidents();

        return (
            <div>
                {incidents}
            </div>
        );
    }
}

class UptimeBar extends React.Component {
    constructor(props) {
        super(props);
    }

    _buildBars() {
        let divs = [];

        const h = this.props.history.slice().reverse();

        let len = 90;

        if (this.props.width < 768)
            len = 30;
        else if (this.props.width < 991)
            len = 60;

        let historyIndex = 0;

        for (let i = 0; i < len; i++) {
            let day = new Date();
            day.setDate(day.getDate() - i);
            day.setHours(0, 0, 0, 0);
                        
            // Subtract days
            let unix = day.getTime();

            let status = 'day-bar-no-data';

            let lastStatus = null;
            let lastTime = null;
            let totalDown = 0;

            while (historyIndex < h.length &&
                h[historyIndex].unix > unix)
            {
                if (!h[historyIndex].up) {
                    status = 'day-bar-down';
                } else if (status != 'day-bar-down') {
                    status = 'day-bar-up';
                }

                if (lastStatus == null) {
                    lastStatus = h[historyIndex].up;
                    if (!lastStatus) {
                        lastTime = h[historyIndex].unix;
                    }
                } else if (!lastStatus && h[historyIndex].up) {
                    lastStatus = true
                    totalDown += lastTime - h[historyIndex].unix;
                } else if (lastStatus && !h[historyIndex].up) {
                    lastStatus = false;
                    lastTime = h[historyIndex].unix;
                }
                
                historyIndex++;
            }

            if (lastStatus === false) {
                totalDown += lastTime - day;
            }

            let title = day.toLocaleDateString();

            if (totalDown) {
                let seconds = totalDown / 1000;
                let text = seconds.toString() + "s";

                if (seconds > 60) {
                    let minutes = seconds / 60;
                    text = minutes.toFixed(0) + "m";

                    if (minutes > 60) {
                        let hours = minutes / 60;
                        text = hours.toFixed(0) + "h";
                    }
                }

                title += "\n\n" + text + " caídos.";
            }

            divs.unshift(
                <div className={"day-bar " + status} title={title} key={i}></div>
            );
        }

        return divs;
    }

    render() {
        if (!this.props.history) {
            return <span></span>
        }

        const bars = this._buildBars();

        return (
            <div className="row">
                <div className="bar-holder">
                    {bars}
                </div>
            </div>
        );
    }
}

class UptimeMetric extends React.Component {
    constructor(props) {
        super(props);
    }

    _calculateUptime() {
        let start = this.props.history[0].unix;
        let end = this.props.history[this.props.history.length - 1].unix;
        let delta = end - start;

        let downtime = 0;

        let lastStatus = this.props.history[0].up;
        let downtimeStart = 0;
        for (let p of this.props.history) {
            if (p.up ^ lastStatus) {
                if (lastStatus) {
                    downtimeStart = p.unix;
                } else {
                    downtime += (p.unix - downtimeStart)
                    downtimeStart = 0;
                }
            }
            lastStatus = p.up;
        }

        if (downtimeStart) downtime += (end - downtimeStart);
        downtime = delta - downtime;

        return Number((100 * downtime / delta).toFixed(3)).toString() + '%';
    }

    render() {
        if (!this.props.history) {
            return <span className="text-muted">Cargando...</span>
        }

        let uptime = this._calculateUptime();

        let len = 90;

        if (this.props.width < 768)
            len = 30;
        else if (this.props.width < 991)
            len = 60;

        return (
            <div className="row">
                <div className="col text-muted">
                    Hace {len} días
                </div>
                <div className="col text-center text-muted">
                    {uptime} uptime
                </div>
                <div className="col text-right text-muted">
                    Hoy
                </div>
            </div>
        );

    }
}

class ServiceItem extends React.Component {
    constructor(props) {
        super(props);
        this.state = { width: 0, height: 0 };
        this.updateWindowDimensions = this.updateWindowDimensions.bind(this);
    }

    render() {
        const statusText = this.props.service.up ? "Operacional" : "Caído";

        let textStatus = "text-success";

        if (!this.props.service.up) {
            textStatus = "text-danger";
        }

        return (
            <li className="list-group-item" key={this.props.service.url}>
                <div className="row">
                    <div className="col" title={this.props.service.url}>
                        {this.props.service.name}
                    </div>
                    <div className={"col text-right " + textStatus}>
                        {statusText}
                    </div>
                </div>
                <UptimeBar history={this.props.history} width={this.state.width} height={this.state.height} onPopup={this.props.onPopup} />
                <UptimeMetric history={this.props.history} width={this.state.width} height={this.state.height} />
            </li>
        );
    }

    componentDidMount() {
        this.updateWindowDimensions();
        window.addEventListener('resize', this.updateWindowDimensions);
    }

    componentWillUnmount() {
        window.removeEventListener('resize', this.updateWindowDimensions);
    }

    updateWindowDimensions() {
        this.setState({ width: window.innerWidth, height: window.innerHeight });
    }
}

class ServiceList extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        let listed = [];

        for (let e in this.props.pages) {
            listed.push(this.props.pages[e]);
        }

        const listItems = listed.map((service) => {
            let history = null;
            if (this.props.history) {
                history = this.props.history[service.url];
            }
            return <ServiceItem service={service} history={history} onPopup={this.props.onPopup} />;
        });

        return (
            <ul className="list-group">
                {listItems}
            </ul>
        );
    }
}

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = { pages: null, history: null, popup: null };
    }

    componentWillMount() {
        getCurrentStatus().then((p) => this.setState({ pages: p }));
        getHistory().then((h) => this.setState({ history: h }));
    }

    _setPopup(popup) {
        this.setState({popup: popup});
    }

    render() {
        if (this.state.pages == null) {
            return (
                <div>
                    <div className="spinner-border" role="status">
                        <span className="sr-only">Cargando...</span>
                    </div>
                    <p>Cargando...</p>
                </div>
            );
        }

        return (
            <div className="container">
                <div className="row text-center">
                    <div className="col">
                        <img className="img-fluid"
                            src="https://www.etsisi.upm.es/sites/default/files/noticias_tic/logosupm.jpg">
                        </img>
                    </div>
                </div>
                <div className="row">
                    <div className="col">
                        <GlobalStatus pages={this.state.pages} />
                    </div>
                </div>
                <div className="row">
                    <div className="col">
                        <ServiceList pages={this.state.pages} history={this.state.history} onPopup={this._setPopup} />
                    </div>
                </div>


                <h2 className="pt-4">Incidencias Pasadas</h2>
                <PastIncidents history={this.state.history} pages={this.state.pages}/>
                <IncidentPopup popup={this.state.popup}/>
            </div>
        );
    }
}

ReactDOM.render(React.createElement(App), domContainer);