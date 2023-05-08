import { BaseLayout } from '@components/ui/layout';
import Link from 'next/link';
import React from 'react';
import { useAdminOrders } from '@components/hooks';
import { Button, Loader } from '@components/ui/common';

export default function AdminOrderScreen() {
    const { data: orders } = useAdminOrders();

    return (
        <>
            <div className="grid md:grid-cols-4 md:gap-5">
                <div>
                    <ul>
                        <li>
                            <Link href="/admin/dashboard">Dashboard</Link>
                        </li>
                        <li>
                            <Link href="/admin/orders">
                                <a className="font-bold">Orders</a>
                            </Link>
                        </li>
                        <li>
                            <Link href="/admin/products">Products</Link>
                        </li>
                        <li>
                            <Link href="/admin/users">Users</Link>
                        </li>
                    </ul>
                </div>
                <div className="overflow-x-auto md:col-span-3">
                    <h1 className="mb-4 text-xl">Admin Orders</h1>

                    {!orders.hasInitialResponse ? (
                        <Loader/>
                    ) : orders.error ? (
                        <div className="alert-error">{orders.error}</div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="min-w-full">
                                <thead className="border-b">
                                    <tr>
                                        <th className="px-5 text-left">ID</th>
                                        <th className="p-5 text-left">USER</th>
                                        <th className="p-5 text-left">DATE</th>
                                        <th className="p-5 text-left">TOTAL</th>
                                        <th className="p-5 text-left">PAID</th>
                                        <th className="p-5 text-left">DELIVERED</th>
                                        <th className="p-5 text-left">ACTION</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {orders.data.map((order) => (
                                        <tr key={order._id} className="border-b">
                                            <td className="p-5">{order._id.substring(20, 24)}</td>
                                            <td className="p-5">
                                                {order.user ? order.user.name : 'DELETED USER'}
                                            </td>
                                            <td className="p-5">
                                                {order.createdAt.substring(0, 10)}
                                            </td>
                                            <td className="p-5">${order.totalPrice}</td>
                                            <td className="p-5">
                                                {order.isPaid
                                                    ? `${order.paidAt.substring(0, 10)}`
                                                    : 'not paid'}
                                            </td>
                                            <td className="p-5">
                                                {order.isDelivered
                                                    ? `${order.deliveredAt.substring(0, 10)}`
                                                    : 'not delivered'}
                                            </td>
                                            <td className="p-5">
                                                <Link href={`/order/${order._id}`} passHref>
                                                    <Button size="sm" variant="white">Details</Button>
                                                </Link>
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

AdminOrderScreen.auth = { adminOnly: true };
AdminOrderScreen.Layout = BaseLayout