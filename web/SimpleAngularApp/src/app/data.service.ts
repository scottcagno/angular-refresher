import {EventEmitter, Injectable} from '@angular/core';
import {Book} from "./model/book";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  books :Array<Book>;
  bookAddedEvent = new EventEmitter<Book>();
  bookDeletedEvent = new EventEmitter<Book>();

  constructor() {
    this.books = new Array<Book>;
    const book1 = new Book('first book', 'matt', 3.99);
    const book2 = new Book('second book', 'james', 5.99);
    const book3 = new Book('third book', 'laura', 8.99);
    this.books.push(book1, book2, book3);
  }

  addBook(book :Book) {
    if (book.author == 'james') {
      this.bookAddedEvent.error('books by james are not allowed!');
    } else {
      this.books.push(book);
      this.bookAddedEvent.emit(book);
    }
  }

  delBook() {
    if (this.books.length > 0) {
      const book = this.books.pop();
      this.bookDeletedEvent.emit(book);
    } else {
      this.bookDeletedEvent.error('there are no more books to remove!');
    }
  }
}
