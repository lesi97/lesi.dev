export function getKillCounts(plugObjectives: any, identifiers: number[]): number {
    // console.log(plugObjectives);
    for (const id of identifiers) {
        const progress = plugObjectives?.[id]?.[0]?.progress;
        if (progress !== undefined) {
            return progress;
        }
    }
    return 0;
}
