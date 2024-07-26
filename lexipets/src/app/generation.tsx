'use client'

import {Suspense, useState} from "react";
import errorEntry from "next/dist/server/typescript/rules/error";

type Pet = {id: string, name: string, species: {id: string, name: string, features: Array<Map<string, string>>}, genes: {feature: Array<Map<string, string>>, dominant: boolean, recessive: boolean}, img: string};

let petName = "";

export function PetComponent() {

    const [pet, setPet] = useState<Pet>();

    function generatePet(e: any) {
        e.preventDefault();

        const form = e.target;
        const formData = new FormData(form);

        let entries = Object.fromEntries(formData.entries())

        fetch('http://localhost:8080/pets/generate', {method: 'POST', body: JSON.stringify({name: entries.petName})}).then(async (res) => {
                if (!res.ok) {
                    throw new Error('Failed to generate pet')
                }
                setPet(await res.json())
            }
        )


    }

    return (
        <div>
            <form method="post" onSubmit={generatePet}>
                <label>
                    "Name your pet"
                    <input name="petName"/>
                </label>
                <button type="submit">Generate your Pet!</button>
                <div className="z-10 w-full max-w-5xl items-center justify-between font-mono text-sm lg:flex">
                    <Suspense fallback={<div>Loading...</div>}>
                        {pet ? pet.name : ""} the {pet ? pet.species.name : ""}:
                        <p>{pet ? pet.img : ""}</p>
                    </Suspense>
                </div>
            </form>
        </div>
)
}