import { ReplaySubject, Observable } from 'rxjs';
import { DataQueryResponse, CircularDataFrame, KeyValue } from '@grafana/data';
import { AWGQuery } from './types';

export class ObservableSubject {
    subject = new ReplaySubject<DataQueryResponse>(1);
    data: KeyValue<CircularDataFrame> = {}; 
    filters: KeyValue<string> = {};
    constructor(request: AWGQuery) {
        this.filters = this.parseQuery(request.queryParam || '');
    }

    parseQuery(queryString: string): KeyValue<string> {
        var query: KeyValue<string> = {};
        var pairs = (queryString[0] === '?' ? queryString.substr(1) : queryString).split('&');
        for (var i = 0; i < pairs.length; i++) {
          var pair: string[] = pairs[i].split('=');
          query[decodeURIComponent(pair[0])] = decodeURIComponent(pair[1] || '');
        }
        return query;
    }

    asObservable(): Observable<DataQueryResponse> {
        return this.subject.asObservable();
    }
}