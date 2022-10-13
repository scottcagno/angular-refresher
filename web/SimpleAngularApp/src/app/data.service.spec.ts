import { TestBed } from '@angular/core/testing';

import { DataService } from './data.service';
import {Book} from "./model/book";

describe('DataService', () => {
  let service: DataService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(DataService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('addBook increases the size of the books array', ()=>{
    const book = new Book();
    service.addBook(book);
    //fail(); // to fail a test
    expect(service.books.length).toEqual(4);
  });
});
