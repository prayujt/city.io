import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';

import { environment } from '../../../environments/environment';
import { CookieService } from 'ngx-cookie-service';
import { Building } from '../../services/building';
import { CityService } from '../../services/city.service';
import { SidebarComponent } from './sidebar/sidebar.component';

@Component({
    selector: 'app-game',
    templateUrl: './game.component.html',
    styleUrls: ['./game.component.css'],
})
export class GameComponent {
    constructor(
        private http: HttpClient,
        private router: Router,
        private cookieService: CookieService,
        private cityService: CityService
    ) {
        this.createCity();
    }
    public row: number = 0;
    public column: number = 0;
    public sessionId: string = '';

    public ngOnInit(): void {
        this.sessionId = this.getID();
        if (this.sessionId != '') {
            this.http
                .get<any>(
                    `http://${environment.API_HOST}:${environment.API_PORT}/sessions/${this.sessionId}`
                )
                .subscribe((response) => {
                    if (!response.status) {
                        this.router.navigate(['login']);
                    }
                });
            this.http
                .get<any>(
                    `http://${environment.API_HOST}:${environment.API_PORT}/cities/${this.sessionId}/buildings`
                )
                .subscribe((response) => {
                    this.cityService.setBuildings(response.buildings);
                });
            setInterval(() => {
                this.http
                    .get<any>(
                        `http://${environment.API_HOST}:${environment.API_PORT}/cities/${this.sessionId}/buildings`
                    )
                    .subscribe((response) => {
                        this.cityService.setBuildings(response.buildings);
                    });
            }, 1000);
        } else {
            this.router.navigate(['login']);
        }
    }

    public getID(): string {
        return this.cookieService.get('cookie');
    }

    createCity(): GameComponent {
        this.cityService.createCity();
        return this;
    }

    get buildings(): Building[][] {
        return this.cityService.getBuildings();
    }

    buildingType!: string;
    buildingLevel!: number;
    buildingProduction!: number;
    happinessChange!: number;
    startTime!: string;
    endTime!: string;
    clicked: boolean = false;
    progBar: boolean = false;

    startUnix!: number;
    endUnix!: number;
    progress!: number;
    interval: any;
    remaining!: number;
    hh!: number;
    mm!: number;
    ss!: number;

    onTileClick(row: number, column: number) {
        // change CSS for selected tile
        const tiles = Array.from(document.querySelectorAll('.tile'));
        tiles.forEach((tile) => tile.classList.remove('selected'));
        const tile = document.querySelector(`td[id="${row} ${column}"]`);
        tile?.classList.add('selected');

        this.row = row;
        this.column = column;
        // update sidebar stats
        // clearInterval(this.interval);
        // this.http
        //     .get<any>(
        //         `http://${environment.API_HOST}:${environment.API_PORT}/cities/${this.ID}/buildings/${row}/${column}`
        //     )
        //     .subscribe((response) => {
        //         this.buildingType = response.buildingType;
        //         this.buildingLevel = response.buildingLevel;
        //         this.buildingProduction = response.buildingProduction;
        //         this.happinessChange = response.happinessChange;
        //         this.startTime = response.startTime;
        //         this.endTime = response.endTime;
        //     });
        // this.clicked = true;

        // this.interval = setInterval(() => {
        //     if (this.startTime != '' && this.endTime != '') {
        //         this.progBar = true;
        //         this.startUnix = Date.parse(this.startTime);
        //         this.endUnix = Date.parse(this.endTime);
        //         let time: number = Date.now();
        //         this.progress =
        //             ((time - this.startUnix) /
        //                 (this.endUnix - this.startUnix)) *
        //             100;
        //         if (this.progress > 100) this.progBar = false;
        //         this.remaining = this.endUnix - time;
        //         this.ss = this.remaining / 1000;
        //         this.hh = Math.floor(this.ss / 3600);
        //         this.ss %= 3600;
        //         this.mm = Math.floor(this.ss / 60);
        //         this.ss %= 60;
        //         this.ss = Math.floor(this.ss);
        //     } else this.progBar = false;
        // }, 100);
    }
}
