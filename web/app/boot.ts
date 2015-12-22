import {bootstrap}    from 'angular2/platform/browser'
import {App} from './app'
import {HTTP_PROVIDERS, JSONP_PROVIDERS} from 'angular2/http';

bootstrap(App, [HTTP_PROVIDERS, JSONP_PROVIDERS]);
