import {Component, OnDestroy, OnInit} from '@angular/core';
import {DataService} from "../data.service";
import {Subscription} from "rxjs";
import {Book} from "../model/book";

@Component({
  selector: 'app-page3',
  templateUrl: './page3.component.html',
  styleUrls: ['./page3.component.css']
})
export class Page3Component implements OnInit, OnDestroy {

  subscription !:Subscription;

  constructor(private dataService: DataService) {}

  ngOnInit(): void {
    this.subscription = this.dataService.bookDeletedEvent.subscribe(
      (book) => {
        console.log(`Book by ${book.author} has been removed`);
      },
      (error) => {
        // do something here...
        console.log(`Got an error: ${error} when trying to delete a book`);
      },
      () => {}
    );
  }

  ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }

  deleteLastBook() {
    this.dataService.delBook();
  }

}
