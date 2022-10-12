import {Component, OnDestroy, OnInit} from '@angular/core';
import {DataService} from "../data.service";
import {Book} from "../model/book";
import {Subscription} from "rxjs";

@Component({
  selector: 'app-page1',
  templateUrl: './page1.component.html',
  styleUrls: ['./page1.component.css']
})
export class Page1Component implements OnInit, OnDestroy {

  pageName = 'Page 1';
  books !:Array<Book>;
  booksByMatt !:number;

  subscription !:Subscription;
  subscription2 !:Subscription;

  // runs when a class is created
  constructor(private dataService :DataService) {
  }

  // runs after a class is instantiated
  ngOnInit(): void {
    setTimeout(()=> { this.pageName = 'First Page' }, 5000)
    this.books = this.dataService.books;
    this.booksByMatt = this.books.filter ((it)=>{return it.author === 'matt'}).length;

    this.subscription = this.dataService.bookAddedEvent.subscribe(
      (newBook) => {
        if (newBook.author === 'matt') {
          this.booksByMatt++
        }
    },
      (error) => {
        // do something here...
        console.log(`Got an error: ${error}`);
      },
      () => {}
    );

    this.subscription = this.dataService.bookDeletedEvent.subscribe(
      (book) => {
        if (book.author === 'matt') {
          this.booksByMatt--;
        }
      },
      (error) => {
      },
      () => {}
    );

  }

  onButtonClick() {
    alert('hello - the date today is ' + new Date());
  }

  ngOnDestroy(): void {
    this.subscription.unsubscribe();
    this.subscription2.unsubscribe();
  }

}
