import { Injectable } from '@angular/core';
import { Constructable } from './constructable';
@Injectable({
  providedIn: 'root'
})
export class ConstructableService {
  constructables: Constructable[] = [];
  types: string[] = ['Apartment', 'Hospital', 'School', 'Supermarket', 'Barracks'];
  constructor() {
    for (let i = 0; i < this.types.length; i++) {
      this.constructables[i] = new Constructable(this.types[i]);
    }
  }
}
