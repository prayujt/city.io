import { Injectable } from '@angular/core';
import { Constructable } from './constructable';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { CookieService } from 'ngx-cookie-service';

interface building {
  buildingType: string;
  buildCost: number;
  buildTime: number;
  buildingProduction: number;
  happinessChange: number;
}

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
      .get<any>(`${environment.API_HOST}/cities/buildings/available`, {
                    headers,
                }).subscribe((response) => {
                  this.setConstructables(response);
              });
    }
  }

  setConstructables(buildings: Array<building>) {
    let i = 0;
    let constructables: Constructable[] = [];
    for (let i = 0; i < buildings.length; i++) {
      constructables[i] = new Constructable(buildings[i]);
    }
    this.constructables = constructables;
  }

  getConstructables(): Constructable[] {
    return this.constructables;
  }
}
