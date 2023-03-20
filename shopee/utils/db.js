import User from '@models/User';
import mongoose from 'mongoose';
import bcryptjs from 'bcryptjs';

const connection = {};

async function connect() {
    if (connection.isConnected) {
        console.log('already connected');
        return;
    }
    if (mongoose.connections.length > 0) {
        connection.isConnected = mongoose.connections[0].readyState;
        if (connection.isConnected === 1) {
            console.log('use previous connection');
            return;
        }
        await mongoose.disconnect();
    }
    const db = await mongoose.connect(process.env.MONGODB_URI);
    console.log('new connection');
    connection.isConnected = db.connections[0].readyState;
}

async function disconnect() {
    if (connection.isConnected) {
        if (process.env.NODE_ENV === 'production') {
            await mongoose.disconnect();
            connection.isConnected = false;
        } else {
            console.log('not disconnected');
        }
    }
}
function convertDocToObj(doc) {
    doc._id = doc._id.toString();
    doc.createdAt = doc.createdAt.toString();
    doc.updatedAt = doc.updatedAt.toString();
    return doc;
}

const createUserMongo = async (name, email, password) => {
    await db.connect();

    const existingUser = await User.findOne({ email: email });
    if (existingUser) {
        await db.disconnect();
        return;
    }

    const newUser = new User({
        name,
        email,
        password: bcryptjs.hashSync(password),
        isAdmin: false,
    });

    const user = await newUser.save();
    await db.disconnect();
    return user._id
}

const deleteUserMongo = async (id) => {
    await db.connect();
    const user = await User.findById(id);
    if (user) {
        if (user.email === 'admin@example.com') {
            return;
        }
        await user.remove();
    }
    await db.disconnect();
}

const updateUserMongo = async (user, name, email, password) => {
    await db.connect();
    const toUpdateUser = await User.findById(user._id);
    toUpdateUser.name = name;
    toUpdateUser.email = email;

    if (password) {
        toUpdateUser.password = bcryptjs.hashSync(password);
    }

    await toUpdateUser.save();
    await db.disconnect();
}

const db = { connect, disconnect, convertDocToObj, createUserMongo, deleteUserMongo, updateUserMongo };
export default db;
