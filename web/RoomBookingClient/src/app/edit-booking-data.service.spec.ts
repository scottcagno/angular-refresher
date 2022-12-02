import { TestBed } from '@angular/core/testing';

import { EditBookingDataService } from './edit-booking-data.service';

describe('EditBookingDataService', () => {
  let service: EditBookingDataService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EditBookingDataService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
