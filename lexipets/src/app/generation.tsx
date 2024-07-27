'use client'

import {Suspense, useState} from "react";
import errorEntry from "next/dist/server/typescript/rules/error";

type Pet = {id: string, name: string, species_name: string, species_features: Array<Map<string, string>>, genes: {feature: Array<Map<string, string>>, dominant: boolean, recessive: boolean}, img: string};

export function PetComponent() {

    const [pet, setPet] = useState<Pet>();
    const [saved, setSaved] = useState<boolean>(false)

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

    function savePet() {
        if (!pet) {
            throw new Error('No pet generated, cannot save')
        }
        console.log(JSON.stringify(pet))
        fetch('http://localhost:8080/pets', {method: 'POST', body: JSON.stringify(pet)}).then(async (res) => {
                if (!res.ok) {
                    throw new Error('Failed to save pet')
                }
                setSaved(true)
                pet.id = await res.json()
                setPet(pet)
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
                        {pet ? pet.name : ""} the {pet ? pet.species_name : ""}:
                        <p>{pet ? pet.img : ""}</p>
                    </Suspense>
                </div>
            </form>
            {saved ? (<p>{pet?.name} is saved! {pet?.id}</p>) : (<button onClick={savePet}>Save your Pet!</button>)}
        </div>
    )
}