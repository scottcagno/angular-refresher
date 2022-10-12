import { Injectable } from '@angular/core';
import {Book} from "./model/book";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  books :Array<Book>;

  constructor() {
    this.books = new Array<Book>;
    const book1 = new Book('first book', 'matt', 3.99);
    const book2 = new Book('second book', 'james', 5.99);
    const book3 = new Book('third book', 'laura', 8.99);
    this.books.push(book1, book2, book3);
  }
}
