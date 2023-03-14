import Link from 'next/link';
import React from 'react';
import { BaseLayout } from '@components/ui/layout';
import { useOrderHistory } from '@components/hooks';
import { Button, Loader, Message } from '@components/ui/common';

function OrderHistoryScreen() {
    const { data: orders } = useOrderHistory();

    return (
        <>
            <h1 className="mb-4 text-xl">Order History</h1>
            {!orders.hasInitialResponse ? (
                <Loader/>
            ) : orders.error ? (
                <Message type="danger">{orders.error}</Message>
            ) : (
                <div className="overflow-x-auto">
                    <table className="min-w-full">
                        <thead className="border-b">
                            <tr>
                                <th className="px-5 text-left">ID</th>
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
                                    <td className=" p-5 ">{order._id.substring(20, 24)}</td>
                                    <td className=" p-5 ">{order.createdAt.substring(0, 10)}</td>
                                    <td className=" p-5 ">${order.totalPrice}</td>
                                    <td className=" p-5 ">
                                        {order.isPaid
                                            ? `${order.paidAt.substring(0, 10)}`
                                            : 'not paid'}
                                    </td>
                                    <td className=" p-5 ">
                                        {order.isDelivered
                                            ? `${order.deliveredAt.substring(0, 10)}`
                                            : 'not delivered'}
                                    </td>
                                    <td className=" p-5 ">
                                        <Link href={`/order/${order._id}`} passHref>
                                            <Button
                                            size="sm"
                                            variant="purple"
                                            className="ml-2"
                                            >
                                            Details
                                            </Button>
                                        </Link>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </>
    );
}

OrderHistoryScreen.auth = true;
OrderHistoryScreen.Layout = BaseLayout
export default OrderHistoryScreen;
