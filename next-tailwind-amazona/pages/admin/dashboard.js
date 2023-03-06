import Link from 'next/link';
import { Bar } from 'react-chartjs-2';

import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend,
} from 'chart.js';
import React from 'react';
import { BaseLayout } from '@components/ui/layout';
import { useSummary } from '@components/hooks';

ChartJS.register(
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend
);

export const options = {
    responsive: true,
    plugins: {
        legend: {
            position: 'top',
        },
    },
};

function AdminDashboardScreen() {
    const { data: summary } = useSummary();

    const data = summary.data ? {
        labels: summary.data.salesData.map((x) => x._id), // 2022/01 2022/03
        datasets: [
            {
                label: 'Sales',
                backgroundColor: 'rgba(162, 222, 208, 1)',
                data: summary.data.salesData.map((x) => x.totalSales),
            },
        ],
    } : null;

    return (
        <>
            <div className="grid md:grid-cols-4 md:gap-5">
                <div>
                    <ul>
                        <li>
                            <Link href="/admin/dashboard">
                                <a className="font-bold">Dashboard</a>
                            </Link>
                        </li>
                        <li>
                            <Link href="/admin/orders">Orders</Link>
                        </li>
                        <li>
                            <Link href="/admin/products">Products</Link>
                        </li>
                        <li>
                            <Link href="/admin/users">Users</Link>
                        </li>
                    </ul>
                </div>
                <div className="md:col-span-3">
                    <h1 className="mb-4 text-xl">Admin Dashboard</h1>
                    {!summary.hasInitialResponse ? (
                        <div>Loading...</div>
                    ) : summary.error ? (
                        <div className="alert-error">{summary.error}</div>
                    ) : (
                        <div>
                            <div className="grid grid-cols-1 md:grid-cols-4">
                                <div className="card m-5 p-5">
                                    <p className="text-3xl">${summary.data.ordersPrice} </p>
                                    <p>Sales</p>
                                    <Link href="/admin/orders">View sales</Link>
                                </div>
                                <div className="card m-5 p-5">
                                    <p className="text-3xl">{summary.data.ordersCount} </p>
                                    <p>Orders</p>
                                    <Link href="/admin/orders">View orders</Link>
                                </div>
                                <div className="card m-5 p-5">
                                    <p className="text-3xl">{summary.data.productsCount} </p>
                                    <p>Products</p>
                                    <Link href="/admin/products">View products</Link>
                                </div>
                                <div className="card m-5 p-5">
                                    <p className="text-3xl">{summary.data.usersCount} </p>
                                    <p>Users</p>
                                    <Link href="/admin/users">View users</Link>
                                </div>
                            </div>
                            <h2 className="text-xl">Sales Report</h2>
                            <Bar
                                options={{
                                    legend: { display: true, position: 'right' },
                                }}
                                data={data}
                            />
                        </div>
                    )}
                </div>
            </div>
        </>
    );
}

AdminDashboardScreen.auth = { adminOnly: true };
AdminDashboardScreen.Layout = BaseLayout
export default AdminDashboardScreen;
