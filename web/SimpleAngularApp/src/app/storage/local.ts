function localSet(key :string, val :any) {
    localStorage.setItem(key, JSON.stringify(val));
}

function localGet(key :string) :any {
  let ret = localStorage.getItem(key);
  if (ret != null) {
    return JSON.parse(ret);
  }
}

function localDel(key :string) {
  localStorage.removeItem(key)
}

export class localDB {

  dbVersion :number = 1;
  dbDefaultStore :string = "default";
  db !:IDBDatabase;

  constructor(name :string) {
    // Let us open our database
    const dbReq = window.indexedDB.open(name, this.dbVersion);
    // Register two event handlers to act on the database being opened successfully, or not
    dbReq.onerror = (e) => {
      console.log(`Error opening database ${e}`);
    };
    dbReq.onsuccess = (e) => {
      console.log(`Database initialised: ${e}`);
      this.db = dbReq.result;

      // register db level even triggers
      this.db.onerror = (e) => { console.log(`Database error: ${e}`) };
      this.db.onabort = (e) => { console.log(`Database abort: ${e}`) };
      this.db.onclose = (e) => { console.log(`Database close: ${e}`) };

      // See if we must create the default store
      if (this.db.objectStoreNames.length > 0) {
        if (!this.db.objectStoreNames.contains(this.dbDefaultStore)) {
          this.db.createObjectStore(this.dbDefaultStore);
        }
      }

    };

  }

  createStore(name :string, pk ?:string) {
    if (pk) {
      this.db.createObjectStore(name, {keyPath:pk});
    } else {
      this.db.createObjectStore(name);
    }
  }

  deleteStore(name :string) {
    this.db.deleteObjectStore(name);
  }

  putData(d :any) {
    // open transaction
    const tx = this.db.transaction([this.dbDefaultStore], 'readwrite');
    tx.oncomplete = (e) => { console.log(`Transaction complete: ${e}`) };
    tx.onerror = (e) => { console.log(`Transaction error: ${e}`) };
    // get our object store and add the data
    const req = tx.objectStore(this.dbDefaultStore).put(d);
    req.onerror = (e) => { console.log(`Error putting data: ${e}`) };
    tx.commit();
  }

  getData(d :any) :any {
    // open transaction
    // open transaction
    const tx = this.db.transaction([this.dbDefaultStore], 'readwrite');
    tx.oncomplete = (e) => { console.log(`Transaction complete: ${e}`) };
    tx.onerror = (e) => { console.log(`Transaction error: ${e}`) };
    // get our object store and attempt to get the data
    const req = tx.objectStore(this.dbDefaultStore).get(d);
    req.onerror = (e) => { console.log(`Error getting data: ${e}`) };
    return req.result;
  }

  delData(d :any) {
    // open transaction
    const tx = this.db.transaction([this.dbDefaultStore], 'readwrite');
    tx.oncomplete = (e) => { console.log(`Transaction complete: ${e}`) };
    tx.onerror = (e) => { console.log(`Transaction error: ${e}`) };
    // get our object store and attempt to get the data
    const req = tx.objectStore(this.dbDefaultStore).delete(d);
    req.onerror = (e) => { console.log(`Error deleting data: ${e}`) };
    tx.commit();
  }

  closeDB() {
    this.db.close();
  }

}
