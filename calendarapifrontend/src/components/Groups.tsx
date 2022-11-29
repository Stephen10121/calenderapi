export type GroupsType = {
    groupId: string;
    groupName: string;
    othersCanAdd: boolean;
}

function Group({ group }: { group: GroupsType }) {
    return(
    <div className="group">
        <h1>{group.groupName}</h1>
        <p>Others {group.othersCanAdd ? "can" : "cannot"} add.</p>
        <p>Id: {group.groupId}</p>
    </div>);
}

export default function Groups({ groups }: {groups: GroupsType[]}) {
    if (!groups) {
        return <h1>No Groups</h1>
    }
  return (
    <div>
        <h1>Groups</h1>
        {groups.map((group: GroupsType) => <Group key={group.groupId} group={group}/>)}
    </div>
  )
}
