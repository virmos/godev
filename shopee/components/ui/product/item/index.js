/* eslint-disable @next/next/no-img-element */
import { Button } from '@components/ui/common';
import Link from 'next/link';
import React from 'react';

export default function ProductItem({ product, addToCartHandler }) {
    return (
        <div className="card group item">
            <Link href={`/product/${product.slug}`}>
                <a>
                    <img
                        src={product.image}
                        alt={product.name}
                        className="rounded shadow object-cover h-64 w-full duration-200 md:block group-hover:scale-110"
                    />
                </a>
            </Link>
            <div className="flex flex-col items-center justify-center p-5 text-center">
                <Link href={`/product/${product.slug}`}>
                    <a>
                        <h2 className="text-lg mb-2 h-14 overflow-hidden">
                            {product.name.substring(0, 38)}...
                        </h2>
                    </a>
                </Link>
                <div className="flex flex-col items-center">
                    <p className="mb-2 text-xl inline-block p-4 py-2 rounded-full font-bold bg-green-200 text-green-700">
                    {product.brand}
                    </p>
                    <p className="block text-indigo-600 lg:inline">
                        ${product.price}
                    </p>
                    <Button
                        onClick={() => addToCartHandler(product)}
                        variant="purple"
                    >
                        Add to cart
                    </Button>
                </div>
            </div>
            {/* <!-- Item Gradient --> */}
            {/* <div className="item-gradient"></div> */}
        </div>
    );
}
