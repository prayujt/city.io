import { Injectable } from '@angular/core';
import { Constructable } from './constructable';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { CookieService } from 'ngx-cookie-service';

@Injectable({
  providedIn: 'root'
})
export class ConstructableService {
  constructables: Constructable[] = [];
  public sessionId: string = '';

  constructor(
    private http: HttpClient,
    private cookieService: CookieService
  ) {
    this.cookieService.delete('cityName');
    let jwtToken = this.cookieService.get('jwtToken');
    if (jwtToken != '') {
      let headers = new HttpHeaders();
      this.http
      .get<any>(`${environment.API_HOST}/cities/buildings`, {
                    headers,
                }).subscribe((response) => {
                  this.setConstructables(response);
              });
    }
  }

  setConstructables(buildings: Map<string, number[]>): Constructable[] {
    let i = 0;
    let constructables: Constructable[] = [];
    buildings.forEach(function(value, key) {
      constructables[i] = new Constructable(key, value);
      i++;
    });
    return constructables;
  }
}
