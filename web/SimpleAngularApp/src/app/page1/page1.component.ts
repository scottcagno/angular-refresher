import { Component, OnInit } from '@angular/core';
import {DataService} from "../data.service";

@Component({
  selector: 'app-page1',
  templateUrl: './page1.component.html',
  styleUrls: ['./page1.component.css']
})
export class Page1Component implements OnInit {

  pageName = 'Page 1';

  // runs when a class is created
  constructor(private dataService :DataService) {
  }

  // runs after a class is instantiated
  ngOnInit(): void {
    setTimeout(()=> { this.pageName = 'First Page' }, 5000)
  }

  onButtonClick() {
    alert('hello - the date today is ' + new Date());
  }

  numberOfBooks() :number {
    return this.dataService.books.length;
  }

}
