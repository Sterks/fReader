db.createUser(
    {
        user: "xml",
        pwd: "rootpassword",
        roles: [
            {
                role: "readWrite",
                db: "readerXML"
            }
        ]
    }
);
