import {bootstrap, Component, NgIf, View} from 'angular2/angular2';
import {HTTP_BINDINGS, Http} from 'angular2/http';

@Component({
    selector: 'div#main',
    bindings: [HTTP_BINDINGS]
})

@View({
    templateUrl: './app.html',
    directives: [NgIf]
})

export class App {
    weather: Object;
    time: Date = new Date();

    constructor(http:Http) {
        this.refreshDate();

        http.get('weather.json').toRx().subscribe(res => {
            this.weather = res.json();
        });
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

bootstrap(App)
    .catch(err => console.error(err));
