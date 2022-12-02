import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {User} from "../../../model/User";
import {Router} from "@angular/router";
import {DataService} from "../../../data.service";

@Component({
  selector: 'app-user-detail',
  templateUrl: './user-detail.component.html',
  styleUrls: ['./user-detail.component.css']
})
export class UserDetailComponent implements OnInit {

  @Input()
  user !:User;

  @Output()
  dataChangedEvent = new EventEmitter();

  message = '';

  constructor(private dataService :DataService, private router :Router) { }

  ngOnInit(): void {
  }

  editUser(){
    this.router.navigate(['admin','users'], {queryParams:{id: this.user.id, action: 'edit'}});
  }

  deleteUser() {
    const result = confirm('Are you sure you wish to delete this user?');
    if (result) {
      this.message = 'deleting user...';
      this.dataService.deleteUser(this.user.id).subscribe(
        next => {
          this.dataChangedEvent.emit();
          this.router.navigate(['admin','users'] );
          this.message = '';
        },
        error => {
          this.message = 'Sorry, this user cannot be deleted at this time.';
        }
      );
    }
  }

  resetPassword() {
    this.message = 'Please wait, resetting password...'
    this.dataService.resetPassword(this.user.id).subscribe(
      next => {
        this.message = `${this.user.name}'s password has been reset`;
        //console.log('password for '+ this.user.name +' has been reset');
        //this.router.navigate(['admin','users']);
      },
      error => {
        this.message = 'Sorry, something went wrong';
      }
    );
  }

}
