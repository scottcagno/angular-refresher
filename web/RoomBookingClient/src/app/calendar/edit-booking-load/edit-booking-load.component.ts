import { Component, OnInit } from '@angular/core';
import {EditBookingDataService} from "../../edit-booking-data.service";
import {ActivatedRoute, Router} from "@angular/router";

@Component({
  selector: 'app-edit-booking-load',
  templateUrl: './edit-booking-load.component.html',
  styleUrls: ['./edit-booking-load.component.css']
})
export class EditBookingLoadComponent implements OnInit {

  stackOverflowProtectionLimit = 10;

  constructor(
    private editBookingDataService : EditBookingDataService,
    private router : Router,
    private route : ActivatedRoute) { }

  ngOnInit(): void {
    setTimeout(()=>{
      this.navigateWhenReady();
    }, 1000)
  }

  navigateWhenReady() {
    if (this.stackOverflowProtectionLimit <= 0) {
      return
    }
    this.stackOverflowProtectionLimit--
    if (this.editBookingDataService.dataLoaded == 2) {
      const id = this.route.snapshot.queryParams['id'];
      if (id) {
        this.router.navigate(['edit-booking'],{queryParams:{id:id}});
      } else {
        this.router.navigate(['add-booking']);
      }
    } else {
      setTimeout(()=>{
        this.navigateWhenReady();
      }, 500);
    }
  }

}
