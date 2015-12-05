import {bootstrap, Component, NgIf, View} from 'angular2/angular2';
import {HTTP_BINDINGS, Http, Response} from 'angular2/http';

@Component({
    selector: 'div#main',
    bindings: [HTTP_BINDINGS]
})

@View({
    templateUrl: './app/app.html',
    directives: [NgIf]
})

export class App {
    weather:    Object;
    departures: Object;
    time:       Date = new Date();

    constructor(http:Http) {
        this.refreshDate();

        http.get('weather.json').map((res: Response) => res.json()).subscribe(res => this.weather = res);
        http.get('vbb.json').map((res: Response) => res.json()).subscribe(res => this.departures = res);
    }

    refreshDate() {
        setTimeout(function() {
            this.time = new Date();
            this.refreshDate();
        }.bind(this), 1000);
    }

    getTime() {
        return ("0" + this.time.getHours()).slice(-2) +
            ":" + ("0" + this.time.getMinutes()).slice(-2) +
            ":" + ("0" + this.time.getSeconds()).slice(-2)
            ;
    }
}

bootstrap(App).catch(err => console.error(err));
