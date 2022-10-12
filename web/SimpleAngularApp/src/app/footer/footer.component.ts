import {Component, Input, OnInit} from '@angular/core';
import {DataService} from "../data.service";
import {Book} from "../model/book";

@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.css']
})
export class FooterComponent implements OnInit {

  @Input()
  lastAccessed!: string;

  constructor(private dataService :DataService) {}

  ngOnInit(): void {
  }

  addBook() {
    const book = new Book('another book', 'matt', 1.99);
    this.dataService.addBook(book);
  }

  addBadBook() {
    const book = new Book('another book', 'james', 1.99);
    this.dataService.addBook(book);
  }

}
