import { BaseLayout } from "@components/ui/layout";
import Link from "next/link";
import { useRouter } from "next/router";
import React, { useCallback, useEffect, useReducer } from "react";
import { toast } from "react-toastify";
import { getError } from "@utils/error";
import { createAdminProducts, deleteAdminProduct } from "@components/api";
import { useAdminProducts } from "@components/hooks";
import { ProductModal } from "@components/ui/order";
import { Button, Loader } from "@components/ui/common";

function reducer(state, action) {
    switch (action.type) {
        case "CREATE_REQUEST":
            return { ...state, loadingCreate: true };
        case "CREATE_SUCCESS":
            return { ...state, loadingCreate: false };
        case "CREATE_FAIL":
            return { ...state, loadingCreate: false };
        case "DELETE_REQUEST":
            return { ...state, loadingDelete: true };
        case "DELETE_SUCCESS":
            return { ...state, loadingDelete: false, successDelete: true };
        case "DELETE_FAIL":
            return { ...state, loadingDelete: false };
        case "DELETE_RESET":
            return { ...state, loadingDelete: false, successDelete: false };

        default:
            state;
    }
}
export default function AdminProductsScreen() {
    const router = useRouter();

    const [{ loadingCreate, successDelete, loadingDelete }, dispatch] =
        useReducer(reducer, {});
    const { data: products } = useAdminProducts();
    const [trigger, setTrigger] = React.useState(false);

    const createHandler = async (product) => {
        try {
            dispatch({ type: "CREATE_REQUEST" });
            const { data } = await createAdminProducts(product);
            dispatch({ type: "CREATE_SUCCESS" });
            toast.success("Product created successfully");
            router.push(`/admin/product/${data.product._id}`);
        } catch (err) {
            dispatch({ type: "CREATE_FAIL" });
            toast.error(getError(err));
        }
    };

    useEffect(() => {
        if (successDelete) {
            dispatch({ type: "DELETE_RESET" });
        }
    }, [successDelete]);

    const deleteHandler = async (productId) => {
        if (!window.confirm("Are you sure?")) {
            return;
        }
        try {
            dispatch({ type: "DELETE_REQUEST" });
            const mutatedProducts = products.data.filter(
                (product) => product._id !== productId
            );
            products.mutate(mutatedProducts);

            await deleteAdminProduct(productId);
            dispatch({ type: "DELETE_SUCCESS" });
            toast.success("Product deleted successfully");
        } catch (err) {
            dispatch({ type: "DELETE_FAIL" });
            toast.error(getError(err));
        }
    };

    const cleanupModal = useCallback(() => {
        setTrigger(false);
    }, []);

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
                            <Link href="/admin/products">
                                <a className="font-bold">Products</a>
                            </Link>
                        </li>
                        <li>
                            <Link href="/admin/users">Users</Link>
                        </li>
                    </ul>
                </div>
                <div className="overflow-x-auto md:col-span-3">
                    <div className="flex justify-between">
                        <h1 className="mb-4 text-xl">Products</h1>
                        {loadingDelete && <div>Deleting item...</div>}
                        <button
                            disabled={loadingCreate}
                            onClick={() => {
                                setTrigger(true);
                            }}
                            className="primary-button"
                        >
                            {loadingCreate ? "Loading" : "Create"}
                        </button>
                    </div>
                    {!products.hasInitialResponse ? (
                        <Loader />
                    ) : products.error ? (
                        <div className="alert-error">{products.error}</div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="min-w-full">
                                <thead className="border-b">
                                    <tr>
                                        <th className="px-5 text-left">ID</th>
                                        <th className="p-5 text-left">NAME</th>
                                        <th className="p-5 text-left">PRICE</th>
                                        <th className="p-5 text-left">
                                            CATEGORY
                                        </th>
                                        <th className="p-5 text-left">COUNT</th>
                                        <th className="p-5 text-left">
                                            RATING
                                        </th>
                                        <th className="p-5 text-left">
                                            ACTIONS
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {products.data.map((product) => (
                                        <tr
                                            key={product._id}
                                            className="border-b"
                                        >
                                            <td className=" p-5 ">
                                                {product._id.substring(20, 24)}
                                            </td>
                                            <td className=" p-5 ">
                                                {product.name}
                                            </td>
                                            <td className=" p-5 ">
                                                ${product.price}
                                            </td>
                                            <td className=" p-5 ">
                                                {product.category}
                                            </td>
                                            <td className=" p-5 ">
                                                {product.countInStock}
                                            </td>
                                            <td className=" p-5 ">
                                                {product.rating}
                                            </td>
                                            <td className=" p-5 ">
                                                <Link
                                                    href={`/admin/product/${product._id}`}
                                                >
                                                    <a>
                                                        <Button
                                                            size="sm"
                                                            variant="white"
                                                        >
                                                            Edit
                                                        </Button>
                                                    </a>
                                                </Link>
                                                &nbsp;
                                                <Button
                                                    onClick={() =>
                                                        deleteHandler(
                                                            product._id
                                                        )
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
            {trigger && (
                <ProductModal
                    trigger={trigger}
                    onSubmit={(formData) => {
                        createHandler(formData);
                        cleanupModal();
                    }}
                    onClose={cleanupModal}
                />
            )}
        </>
    );
}

AdminProductsScreen.auth = { adminOnly: true };
AdminProductsScreen.Layout = BaseLayout;
