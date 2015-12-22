import {Component} from 'angular2/core';
import {Http, Jsonp} from 'angular2/http';

@Component({
    selector: 'div#main',
    templateUrl: './app/app.html'
})

export class App {
    time:       Date = new Date();
    weather:    Object;
    departures: Object;
    quote:      Object;

    constructor(http: Http, jsonp: Jsonp) {
        // time
        this.refreshDate();

        // weather
        http.get('weather.json').subscribe(res => this.weather = res.json());

        // departures
        http.get('vbb.json').subscribe(res => this.departures = res.json());

        // quote
        var quoteApiEndpoint: string = "http://api.forismatic.com/api/1.0/?method=getQuote&lang=en&format=jsonp&jsonp=JSONP_CALLBACK";
        jsonp.get(quoteApiEndpoint).subscribe(res => this.quote = res.json());
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
            ":" + ("0" + this.time.getSeconds()).slice(-2);
    }
}
