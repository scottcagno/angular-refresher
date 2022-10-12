import { Component, OnInit } from '@angular/core';
import {DataService} from "../data.service";
import {Book} from "../model/book";

@Component({
  selector: 'app-page1',
  templateUrl: './page1.component.html',
  styleUrls: ['./page1.component.css']
})
export class Page1Component implements OnInit {

  pageName = 'Page 1';
  books !:Array<Book>;
  booksByMatt !:number;

  // runs when a class is created
  constructor(private dataService :DataService) {
  }

  // runs after a class is instantiated
  ngOnInit(): void {
    setTimeout(()=> { this.pageName = 'First Page' }, 5000)
    this.books = this.dataService.books;
    this.booksByMatt = this.books.filter ((it)=>{return it.author === 'matt'}).length;
    this.dataService.bookAddedEvent.subscribe(
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
  }

  onButtonClick() {
    alert('hello - the date today is ' + new Date());
  }

}
