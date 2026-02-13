import { describe, it, expect, vi } from 'vitest'
import { members as rawMembers, roles, org } from '../assets/about.json'
import * as fs from 'fs'
import * as path from 'path'

// Process members the same way Footer.vue does
const members = rawMembers.map(m => ({
  name: m.name,
  firstName: m.name.split(' ')[0],
  org: org[m.org_id],
  roles: m.role_ids.map(r => roles[r])
}))

describe('About / Team Members', () => {
  it('should include Geoffrey in the members list', () => {
    const geoff = members.find(m => m.firstName === 'Geoffrey')
    expect(geoff).toBeDefined()
    expect(geoff!.name).toContain('Geoffrey')
  })

  it('should have Geoffrey.jpg in the public/about directory', () => {
    const publicAboutDir = path.resolve(__dirname, '../../public/about');
    const files = fs.readdirSync(publicAboutDir);
    expect(files).toContain('Geoffrey.jpg')
  })

  it('should have a matching image for every team member', () => {
    const aboutDir = path.resolve(__dirname, '../../public/about')
    const files = fs.readdirSync(aboutDir)

    for (const member of members) {
      const expectedFile = `${member.firstName}.jpg`
      expect(files, `Missing image for ${member.name}: ${expectedFile}`).toContain(expectedFile)
    }
  })

  it('should have 14 team members total', () => {
    expect(members.length).toBe(14)
  })

  it('Footer.vue should render all members with circular avatar markup', () => {
    const footerSrc = fs.readFileSync(
      path.resolve(__dirname, '../components/Footer.vue'),
      'utf-8'
    )

    // Verify the v-for loop and circular avatar classes exist
    expect(footerSrc).toContain('v-for="(m, i) in members"')
    expect(footerSrc).toContain('rounded-full')
    expect(footerSrc).toContain('border-rose-300')
    expect(footerSrc).toContain("'/about/' + m.firstName + '.jpg'")
  })
})
