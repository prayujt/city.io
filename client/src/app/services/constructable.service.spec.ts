import { TestBed } from '@angular/core/testing';

import { ConstructableService } from './constructable.service';

describe('ConstructableService', () => {
  let service: ConstructableService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ConstructableService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
