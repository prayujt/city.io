import { Component, Input } from '@angular/core';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { CookieService } from 'ngx-cookie-service';

import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/environment';

@Component({
    selector: 'game-sidebar',
    templateUrl: './sidebar.component.html',
    styleUrls: ['./sidebar.component.css'],
})
export class SidebarComponent {
    @Input() row!: number;
    @Input() column!: number;
    @Input() sessionId!: string;

    constructor(
        private http: HttpClient,
        private router: Router,
        private cookieService: CookieService,
        private _snackBar: MatSnackBar
    ) {}

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

    public ngOnInit(): void {
        this.interval = setInterval(() => {
            if (this.buildingType != '') this.clicked = true;
            this.http
                .get<any>(
                    `http://${environment.API_HOST}:${environment.API_PORT}/cities/${this.sessionId}/buildings/${this.row}/${this.column}`
                )
                .subscribe((response) => {
                    this.buildingType = response.buildingType;
                    this.buildingLevel = response.buildingLevel;
                    this.buildingProduction = response.buildingProduction;
                    this.happinessChange = response.happinessChange;
                    this.startTime = response.startTime;
                    this.endTime = response.endTime;
                });

            if (this.startTime != '' && this.endTime != '') {
                this.progBar = true;
                this.startUnix = Date.parse(this.startTime);
                this.endUnix = Date.parse(this.endTime);
                let time: number = Date.now();
                this.progress =
                    ((time - this.startUnix) /
                        (this.endUnix - this.startUnix)) *
                    100;
                if (this.progress > 100) this.progBar = false;
                this.remaining = this.endUnix - time;
                this.ss = this.remaining / 1000;
                this.hh = Math.floor(this.ss / 3600);
                this.ss %= 3600;
                this.mm = Math.floor(this.ss / 60);
                this.ss %= 60;
                this.ss = Math.floor(this.ss);
            } else this.progBar = false;
        }, 100);
    }

    public logOut(): void {
        this.http
            .post<any>(
                `http://${environment.API_HOST}:${environment.API_PORT}/sessions/logout`,
                {
                    sessionId: this.sessionId,
                }
            )
            .subscribe((response) => {
                if (response.status) {
                    this._snackBar.open('Log out successful!', 'Close', {
                        duration: 2000,
                    });
                    this.router.navigate(['login']).then(() => {
                        window.location.reload();
                        this.cookieService.delete('cookie');
                    });
                    this.cookieService.delete('cookie');
                } else {
                    this._snackBar.open('Could not log out!', 'Close', {
                        duration: 2000,
                    });
                }
            });
    }
}
