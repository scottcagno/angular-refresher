import {environment} from "../../environments/environment";

export class User {
  id !:number;
  name !:string;

  constructor(id ?:number, name ?:string) {
    if (id) { this.id = id }
    if (name) { this.name = name }
  }

  getRole():string {
    return 'standard';
  }

  static fromHttp(u: User) :User {
    return new User(u.id, u.name);
  }

  static endpoint(id ?:number):string {
    if (id) {
      return environment.restUrl + `/api/users?id=${id}`
    }
    return environment.restUrl + `/api/users`
  }
}
