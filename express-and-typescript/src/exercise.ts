export type Exercise = {
    name: string;
    alternative?: string;
}

export function genExercises(): Exercise[] {
    const exercises: Exercise[] = [
        {name: 'Jumping Jacks'},
        {name: 'Wall Sit'},
        {name: 'Pushups'},
        {name: 'Crunches'},
        {name: 'Step Up', alternative:'Lateral Lunge'},
        {name: 'Squat'},
        {name: 'Triceps Dip', alternative:'Reverse Tabletop'},
        {name: 'Plank'},
        {name: 'High Knees'},
        {name: 'Lunge'},
        {name: 'T Pushup'},
        {name: 'Side Plank (left side)'},
        {name: 'Side Plank (right side)'}
    ]

    return exercises
}