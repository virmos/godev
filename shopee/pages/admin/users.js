import { BaseLayout } from "@components/ui/layout";
import Link from "next/link";
import React, { useEffect, useReducer } from "react";
import { toast } from "react-toastify";
import { getError } from "@utils/error";
import { useAdminUsers } from "@components/hooks";
import { deleteAdminUser } from "@components/api";
import { Button, Loader } from "@components/ui/common";

function reducer(state, action) {
    switch (action.type) {
        case "DELETE_REQUEST":
            return { ...state, loadingDelete: true };
        case "DELETE_SUCCESS":
            return { ...state, loadingDelete: false, successDelete: true };
        case "DELETE_FAIL":
            return { ...state, loadingDelete: false };
        case "DELETE_RESET":
            return { ...state, loadingDelete: false, successDelete: false };
        default:
            return state;
    }
}

function AdminUsersScreen() {
    const [{ successDelete, loadingDelete }, dispatch] = useReducer(reducer, {
        loading: true,
        users: [],
        error: "",
    });

    const { data: users } = useAdminUsers();

    useEffect(() => {
        if (successDelete) {
            dispatch({ type: "DELETE_RESET" });
        }
    }, [successDelete]);

    const deleteHandler = async (userId) => {
        if (!window.confirm("Are you sure?")) {
            return;
        }
        try {
            dispatch({ type: "DELETE_REQUEST" });
            const mutatedUser = users.data?.filter((user) => {
                user._id !== userId;
            });
            await deleteAdminUser(userId);
            users.mutate(mutatedUser);
            dispatch({ type: "DELETE_SUCCESS" });
            toast.success("User deleted successfully");
        } catch (err) {
            dispatch({ type: "DELETE_FAIL" });
            toast.error(getError(err));
        }
    };

    return (
        <>
            <div className="grid md:grid-cols-4 md:gap-5">
                <div>
                    <ul>
                        <li>
                            <Link href="/admin/dashboard">Dashboard</Link>
                        </li>
                        <li>
                            <Link href="/admin/orders">Orders</Link>
                        </li>
                        <li>
                            <Link href="/admin/products">Products</Link>
                        </li>
                        <li>
                            <Link href="/admin/users">
                                <a className="font-bold">Users</a>
                            </Link>
                        </li>
                    </ul>
                </div>
                <div className="overflow-x-auto md:col-span-3">
                    <h1 className="mb-4 text-xl">Users</h1>
                    {loadingDelete && <div>Deleting...</div>}
                    {!users.hasInitialResponse ? (
                        <Loader />
                    ) : users.error ? (
                        <div className="alert-error">{users.error}</div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="min-w-full">
                                <thead className="border-b">
                                    <tr>
                                        <th className="px-5 text-left">ID</th>
                                        <th className="p-5 text-left">NAME</th>
                                        <th className="p-5 text-left">EMAIL</th>
                                        <th className="p-5 text-left">ADMIN</th>
                                        <th className="p-5 text-left">
                                            ACTIONS
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {users.data?.map((user) => (
                                        <tr key={user._id} className="border-b">
                                            <td className=" p-5 ">
                                                {user._id}
                                            </td>
                                            <td className=" p-5 ">
                                                {user.name}
                                            </td>
                                            <td className=" p-5 ">
                                                {user.email}
                                            </td>
                                            <td className=" p-5 ">
                                                {user.isAdmin ? "YES" : "NO"}
                                            </td>
                                            <td className=" p-5 ">
                                                &nbsp;
                                                <Button
                                                    onClick={() =>
                                                        deleteHandler(user._id)
                                                    }
                                                    size="sm"
                                                    variant="white"
                                                >
                                                    Delete
                                                </Button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                </div>
            </div>
        </>
    );
}

AdminUsersScreen.auth = { adminOnly: true };
AdminUsersScreen.Layout = BaseLayout;
export default AdminUsersScreen;
