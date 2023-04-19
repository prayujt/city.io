import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { BrowserModule, Title } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FlexLayoutModule } from '@angular/flex-layout';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';
import { AppRoutingModule } from './app-routing.module';
import { RouterModule } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { AppComponent } from './app.component';
import { LoginComponent } from './components/login/login.component';
import { RegisterComponent } from './components/register/register.component';
import { GameComponent } from './components/game/game.component';
import {
    AttackDialogComponent,
    CityNameChangeDialogComponent,
    SidebarComponent,
    TrainDialogComponent,
    VisitDialogComponent,
    MarchesDialogComponent,
    DeleteDialogComponent,
    BattleLogsDialogComponent,
    TerritoryDialogComponent,
} from './components/game/sidebar/sidebar.component';

import { AngularMaterialModule } from './angular-material.module';

import { CookieService } from 'ngx-cookie-service';

@NgModule({
    declarations: [
        AppComponent,
        LoginComponent,
        RegisterComponent,
        GameComponent,
        SidebarComponent,
        VisitDialogComponent,
        CityNameChangeDialogComponent,
        AttackDialogComponent,
        MarchesDialogComponent,
        TrainDialogComponent,
        DeleteDialogComponent,
        BattleLogsDialogComponent,
        TerritoryDialogComponent,
    ],
    imports: [
        BrowserModule,
        AppRoutingModule,
        RouterModule,
        RouterTestingModule,
        BrowserAnimationsModule,
        AngularMaterialModule,
        FlexLayoutModule,
        FormsModule,
        ReactiveFormsModule,
        HttpClientModule,
    ],
    providers: [
        CookieService,
        Title,
        RouterModule,
        HttpClientModule,
        MatSnackBar,
    ],
    bootstrap: [AppComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class AppModule {}
