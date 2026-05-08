export namespace config {
	
	export class DownloadDir {
	    path: string;
	    isDefault: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DownloadDir(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.isDefault = source["isDefault"];
	    }
	}

}

export namespace providers {
	
	export class Photo {
	    url: string;
	    thumbnailUrl?: string;
	
	    static createFrom(source: any = {}) {
	        return new Photo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.thumbnailUrl = source["thumbnailUrl"];
	    }
	}

}

