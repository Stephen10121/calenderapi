export type PendingGroupsType = {
    groupId: string;
    groupName: string;
}

function PendingGroup({ group }: { group: PendingGroupsType }) {
    return(
    <div className="group">
        <h1>{group.groupName}</h1>
        <p>Id: {group.groupId}</p>
    </div>);
}

export default function PendingGroups({ groups }: {groups: PendingGroupsType[]}) {
    if (!groups) {
        return <h1>No Pending Groups</h1>
    }
  return (
    <div>
        <h1>Pending Groups</h1>
        {groups.map((group: PendingGroupsType) => <PendingGroup key={group.groupId} group={group}/>)}
    </div>
  )
}
