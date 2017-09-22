#!/usr/bin/env bash

mongo mongodb://localhost:27017/datasets <<EOF
db.dropDatabase();
db.datasets.insert({
    "_id": "95c4669b-3ae9-4ba7-b690-87e890a1c67c",
    "current": {
        "collection_id": "95c4669b-3ae9-4ba7-b690-87e890a1c543",
        "contact": {
            "email": "jsinclair@test.co.uk",
            "name": "john sinclair",
            "telephone": "01633 123456"
        },
        "description": "census covers the ethnicity of people living in the uk",
        "links": {
            "editions": {
                "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions"
            },
            "latest_version": {
                "id": "63294ed7-dccf-4f30-ad57-62365f038fb7",
                "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2016/versions/1"
            },
            "self": {
                "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c"
            }
        },
        "next_release": "2017-08-23",
        "periodicity": "yearly",
        "publisher": {
            "name": "The office of national statistics",
            "type": "goverment department",
            "href": "https://www.ons.gov.uk/"
        },
        "state": "published",
        "theme": "population",
        "title": "CPI",
        "last_updated": "2017-08-23T15:09:11.829+01:00"
    },
    "next": {
        "contact": {
            "email": "jsinclair@test.co.uk",
            "name": "john sinclair",
            "telephone": "01633 123456"
        },
        "description": "census covers the ethnicity of people living in the uk",
        "links": {
            "editions": {
                "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions"
            },
            "latest_version": {
                "id": "679ebe5f-d9cd-4d6e-8afc-6a2a4f991ccf",
                "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2017/versions/1"
            },
            "self": {
                "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c"
            }
        },
        "next_release": "2018-08-23",
        "periodicity": "yearly",
        "publisher": {
            "name": "The office of national statistics",
            "type": "goverment department",
            "href": "https://www.ons.gov.uk/"
        },
        "state": "created",
        "theme": "population",
        "title": "CPI",
        "last_updated": "2017-08-25T15:09:11.829+01:00"
    }
});
db.datasets.insert({
    "_id": "a9fa845c-0c05-4954-aed7-752b8208da34",
    "next": {
        "contact": {
            "email": "ldavis@test.co.uk",
            "name": "lawrence davis",
            "telephone": "01633 123457"
        },
        "description": "census covers the ethnicity of people living in the uk",
        "links": {
            "editions": {
                "href": "http://localhost:22000/datasets/456/editions"
            },
            "latest_version": {
                "id": "3b45921b-0efa-4844-a16e-3f9f30df4f88",
                "href": "http://localhost:22000/datasets/a9fa845c-0c05-4954-aed7-752b8208da34/editions/2011/versions/1"
            },
            "self": {
                "href": "http://localhost:22000/datasets/a9fa845c-0c05-4954-aed7-752b8208da34"
            }
        },
        "next_release": "2021-04-30",
        "periodicity": "yearly",
        "publisher": {
            "name": "The office of national statistics",
            "type": "goverment department",
            "href": "https://www.ons.gov.uk/"
        },
        "theme": "population",
        "title": "CensusEthnicity",
        "state": "created",
        "last_updated": "2017-08-25T15:09:11.829+01:00"
    }
});
db.editions.insert({
    "edition": "2016",
    "_id": "a051a058-58a9-4ba4-8374-fbb7315d3b78",
    "links": {
        "dataset": {
            "id": "95c4669b-3ae9-4ba7-b690-87e890a1c67c",
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c"
        },
        "self": {
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2016"
        },
        "versions": {
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2016/versions"
        }
    },
    "state": "published",
    "last_updated": "2017-08-25T15:09:11.829+01:00"
});
db.editions.insert({
    "edition": "2017",
    "_id": "8af20615-c4c5-4bb9-af35-a2530e5a2433",
    "links": {
        "dataset": {
            "id": "95c4669b-3ae9-4ba7-b690-87e890a1c67c",
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c"
        },
        "self": {
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2017"
        },
        "versions": {
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2017/versions"
        }
    },
    "state": "created",
    "last_updated": "2017-08-25T15:09:11.829+01:00"
});
db.editions.insert({
    "_id": "2dc3a321-2c31-4a8a-9a8d-7962d7590ed3",
    "edition": "2011",
    "links": {
        "dataset": {
            "id": "a9fa845c-0c05-4954-aed7-752b8208da34",
            "href": "http://localhost:22000/datasets/456"
        },
        "self": {
            "href": "http://localhost:22000/datasets/a9fa845c-0c05-4954-aed7-752b8208da34/editions/2011"
        },
        "versions": {
            "href": "http://localhost:22000/datasets/a9fa845c-0c05-4954-aed7-752b8208da34/editions/2011/versions"
        }
    },
    "state": "created",
    "last_updated": "2017-08-25T15:09:11.829+01:00"
});
db.versions.insert({
    "_id": "63294ed7-dccf-4f30-ad57-62365f038fb7",
    "edition": "2016",
    "instance_id": "63294ed7-dccf-4f30-ad57-62365f038111",
    "collection_id": "95c4669b-3ae9-4ba7-b690-87e890a1c543",
    "license": "ONS",
    "links": {
        "dataset": {
            "id": "95c4669b-3ae9-4ba7-b690-87e890a1c67c",
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c"
        },
        "edition": {
            "id": "a051a058-58a9-4ba4-8374-fbb7315d3b78",
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2016"
        },
        "self": {
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2016/versions/1"
        },
        "dimensions": {
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2016/versions/1/dimensions"
        }
    },
    "release_date": "2016-08-23",
    "state": "published",
    "last_updated": "2017-08-25T15:09:11.829+01:00",
    "version": "1"
});
db.versions.insert({
    "_id": "4ce2ee5c-d50b-469f-b005-ddfed3f5072b",
    "edition": "2016",
    "instance_id": "4ce2ee5c-d50b-469f-b005-ddfed3f50111",
    "license": "ONS",
    "links": {
        "dataset": {
            "id": "95c4669b-3ae9-4ba7-b690-87e890a1c67c",
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c"
        },
        "edition": {
            "id": "a051a058-58a9-4ba4-8374-fbb7315d3b78",
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2016"
        },
        "self": {
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2016/versions/2"
        },
        "dimensions": {
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2016/versions/2/dimensions"
        }
    },
    "release_date": "2016-08-24",
    "state": "created",
    "last_updated": "2017-08-25T15:09:11.829+01:00",
    "version": "2"
});
db.versions.insert({
    "_id": "679ebe5f-d9cd-4d6e-8afc-6a2a4f991ccf",
    "edition": "2017",
    "instance_id": "679ebe5f-d9cd-4d6e-8afc-6a2a4f991111",
    "license": "ONS",
    "links": {
        "dataset": {
            "id": "95c4669b-3ae9-4ba7-b690-87e890a1c67c",
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c"
        },
        "edition": {
            "id": "8af20615-c4c5-4bb9-af35-a2530e5a2433",
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2017"
        },
        "self": {
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2017/versions/1"
        },
        "dimensions": {
            "href": "http://localhost:22000/datasets/95c4669b-3ae9-4ba7-b690-87e890a1c67c/editions/2017/versions/1/dimensions"
        }
    },
    "release_date": "2017-08-23",
    "state": "created",
    "last_updated": "2017-08-25T15:09:11.829+01:00",
    "version": "1"
});
db.versions.insert({
    "_id": "3b45921b-0efa-4844-a16e-3f9f30df4f88",
    "edition": "2011",
    "instance_id": "3b45921b-0efa-4844-a16e-3f9f30df4111",
    "license": "ONS",
    "links": {
        "dataset": {
            "id": "a9fa845c-0c05-4954-aed7-752b8208da34",
            "href": "http://localhost:22000/datasets/a9fa845c-0c05-4954-aed7-752b8208da34"
        },
        "edition": {
            "id": "2dc3a321-2c31-4a8a-9a8d-7962d7590ed3",
            "href": "http://localhost:22000/datasets/a9fa845c-0c05-4954-aed7-752b8208da34/editions/2011"
        },
        "self": {
            "href": "http://localhost:22000/datasets/a9fa845c-0c05-4954-aed7-752b8208da34/editions/2011/versions/1"
        },
        "dimensions": {
            "href": "http://localhost:22000/datasets/a9fa845c-0c05-4954-aed7-752b8208da34/editions/2011/versions/1/dimensions"
        }
    },
    "release_date": "2011-04-30",
    "state": "created",
    "last_updated": "2017-08-25T15:09:11.829+01:00",
    "version": "1"
});
db.contacts.insert({
    "_id": "1",
    "name": "john sinclair",
    "email": "jsinclair@test.co.uk",
    "telephone": "01633 123456",
    "last_updated": "2017-08-25T15:09:11.829+01:00"
});
db.contacts.insert({
    "_id": "2",
    "name": "lawrence davis",
    "email": "ldavis@test.co.uk",
    "telephone": "01633 123457",
    "last_updated": "2017-08-25T15:09:11.829+01:00"
});
db.editions.ensureIndex({"links.dataset.id":1},{"background":true});
db.editions.ensureIndex({"edition":1, "links.dataset.id":1},{"background":true});
db.versions.ensureIndex({"links.dataset.id":1},{"background":true});
db.versions.ensureIndex({"edition":1,"links.dataset.id":1},{"background":true});
db.versions.ensureIndex({"version":1,"edition":1,"links.dataset.id":1},{"background":true});
db.versions.ensureIndex({"version":1,"links.edition.id":1},{"background":true})
db.instances.insert({
    "_id": "AB3BAE9B-5C4D-4640-8936-8502D0DB954D",
    "id": "3b45921b-0efa-4844-a16e-3f9f30df4f88",
    "links": {"job": {"id": "260EDB0F-2BCD-4006-B441-571F504273E0", "link": "http://localhost:22000/jobs/260EDB0F-2BCD-4006-B441-571F504273E0"},},
    "state": "created",
    "total_observations": 0,
    "total_inserted_observations": 0,
    "headers": ["V4_1", "time", "age"],
    "telephone": "01633 123457",
    "last_updated": "2017-08-25T15:09:11.829+01:00",
});
db.instances.insert({
    "_id": "0F06AB0E-A5D1-409A-8183-BACDF2326205",
    "id": "665BEE8A-D88F-448C-BF29-186D18B8DABE",
    "links": {"job": {"id": "260EDB0F-2BCD-4006-B441-571F504273E0", "href": "http://localhost:22000/jobs/260EDB0F-2BCD-4006-B441-571F504273E0"}},
    "state": "created",
    "total_observations": 0,
    "total_inserted_observations": 0,
    "headers": ["V4_1", "time", "age"],
    "telephone": "01633 123457",
    "last_updated": "2017-08-25T15:09:11.829+01:00",
});
db.dimensions.insert({
     "_id":"0F06AB0E-A5D1-409A-8183-BACDF2326205",
     "instance_id": "665BEE8A-D88F-448C-BF29-186D18B8DABE",
     "name": "age",
     "value": "32",
     "node_id":"80",
});
db.dimensions.insert({
     "_id":"AB3BAE9B-5C4D-4640-8936-8502D0DB954D",
     "instance_id": "665BEE8A-D88F-448C-BF29-186D18B8DABE",
     "name": "age",
     "value": "37",
     "node_id":"88",
     "last_updated": "2017-08-25T15:09:11.829+01:00",
});
EOF