print("\n ############ Init script started ############ \n")
const database = "authentication"
var uuid = UUID()
    .toString('hex')
    .replace(/^(.{8})(.{4})(.{4})(.{4})(.{12})$/, '$1-$2-$3-$4-$5')

const user = {
    _id: uuid,
    createAt: new Date(),
    updateAt: new Date(),
    username: "admin",
    password: "$2a$12$8x8yeXV1D1RR.gC8J/6Z6.5WvsalnR09jhcYCx7HTtdf5N.oeq3MK",
    roles: ["admin"],
}

db = db.getSiblingDB(database);
print("\n ----- Authentication database created ----- \n")

db.createCollection('user');
print("\n ----- User collection created ----- \n")

db.user.createIndex( { "username": 1 }, { unique: true } )
print("\n ----- Username unique index created ----- \n")

db.user.insertOne(user);
print("\n ----- Admin user inserted ----- \n")

print("\n ############ Init script finished ############ \n")